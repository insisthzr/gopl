//a simple ftpd
//non error handler
//only a few commands
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Connection struct {
	rw       net.Conn
	prevCmd  string
	dataAddr string
	dir      string
}

func (c *Connection) writeln(args ...interface{}) error {
	args = append(args, "\r\n")
	//i write the code before like this:
	// _, err := fmt.Fprint(c.rw, args)
	//fucking kill me
	_, err := fmt.Fprint(c.rw, args...)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (c *Connection) quit() {
	c.writeln("221 Goodbye.")
}

func (c *Connection) user() {
	c.writeln("230 Login successful.")
}

func (c *Connection) port(args []string) {
	if len(args) != 1 {
		c.writeln("501 Usage: PORT a,b,c,d,p1,p2")
		return
	}
	addr, err := getAddrFromFTP(args[0])
	if err != nil {
		c.writeln("501 %v", err)
		return
	}
	c.dataAddr = addr
	c.writeln("200 PORT command successful.")
}

func (c *Connection) list(args []string) {
	if len(args) > 1 {
		c.writeln("501 ")
		return
	}
	filename := c.dir
	if len(args) == 1 {
		filename = path.Join(filename, args[0])
	}
	wrc, err := c.dataConn()
	if err != nil {
		c.writeln("425 Can't open data connection.")
	}
	defer wrc.Close()

	cmd := exec.Command("ls", "-l", filename)
	out := &bytes.Buffer{}
	cmd.Stdout = out
	err = cmd.Run()
	if err != nil {
		log.Println(err)
	}
	c.writeln("150 Here comes the directory listing.")
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		fmt.Fprintf(wrc, "%s\r\n", line) //write to data socket not cmd socket
	}
	c.writeln("226 Closing data connection. List successful.")
}

func (c *Connection) cwd(args []string) {
	if len(args) > 1 {
		c.writeln("501 ")
	}
	dir := c.dir
	if len(args) == 1 {
		dir = path.Join(dir, args[0])
	}
	c.dir = dir
	c.writeln("250 Directory successfully changed.")
}

//LF -> CRLF
func (c *Connection) retr(args []string) {
	if len(args) != 1 {
		c.writeln("501 error")
		return
	}
	rwc, err := c.dataConn()
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}
	defer rwc.Close()
	filePath := path.Join(c.dir, args[0])
	file, err := os.Open(filePath)
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}
	defer file.Close()
	c.writeln("150 File ok. Sending.")
	_, err = io.Copy(rwc, file)
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}
	c.writeln("226 Transfer complete.")
}

func (c *Connection) stor(args []string) {
	if len(args) != 1 {
		c.writeln("501 error")
		return
	}

	rwc, err := c.dataConn()
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}
	defer rwc.Close()

	filePath := path.Join(c.dir, args[0])
	file, err := os.Create(filePath)
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}
	defer file.Close()

	c.writeln("150 Ok to send data.")
	_, err = io.Copy(file, rwc)
	if err != nil {
		c.writeln("501 %s", err.Error())
		return
	}

	c.writeln("226 Transfer complete.")
}

func (c *Connection) dataConn() (io.ReadWriteCloser, error) {
	var conn io.ReadWriteCloser
	var err error
	switch c.prevCmd {
	case "PORT":
		conn, err = net.Dial("tcp4", c.dataAddr)
		if err != nil {
			log.Println(conn)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("previous command not PORT")
	}
	return conn, nil
}

func (c *Connection) run() {
	log.Println("client connected")
	c.writeln("200 Ready.")
	scanner := bufio.NewScanner(c.rw)
	for scanner.Scan() {
		text := scanner.Text()
		log.Println(text)
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}
		cmd := strings.ToUpper(fields[0])
		args := []string{}
		if len(fields) > 1 {
			args = fields[1:]
		}
		switch cmd {
		case "QUIT":
			c.quit()
			break
		case "USER":
			c.user()
		case "PORT":
			c.port(args)
		case "LIST":
			c.list(args)
		case "CWD":
			c.cwd(args)
		case "RETR":
			c.retr(args)
		case "STOR":
			c.stor(args)
		default:
			c.writeln(fmt.Sprintf("502 Command %q not implemented.", cmd))
		}
		c.prevCmd = cmd
	}
	log.Println("client closed")
}

func NewConnection(conn net.Conn) *Connection {
	dir, err := os.Getwd()
	if err != nil {
		dir = "/"
	}
	return &Connection{rw: conn, dir: dir}
}

func getAddrFromFTP(address string) (string, error) {
	var a, b, c, d byte
	var p1, p2 int
	_, err := fmt.Sscanf(address, "%d,%d,%d,%d,%d,%d", &a, &b, &c, &d, &p1, &p2)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d.%d.%d.%d:%d", a, b, c, d, 256*p1+p2), nil
}

var (
	port = flag.String("port", "8000", "port")
)

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp4", ":"+*port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("listening at %s\n", *port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		c := NewConnection(conn)
		go c.run()
	}
}
