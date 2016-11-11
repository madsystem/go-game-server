package main
import ("bufio" 
        "net")

type Client struct {
	chanInAction chan string 
	chanOutAction chan string 
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func (client *Client) Read() {
	for {
		line, error := client.reader.ReadString('\r')
        if error != nil {
            // todo close + remove client
            break
        }
		client.chanInAction <- line
	}
}

func (client *Client) Write() {
	for data := range client.chanOutAction {
		_, error := client.writer.WriteString(data)
        if error != nil {
            // todo close + remove client
            break
        }
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}


func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		chanInAction: make(chan string),
		chanOutAction: make(chan string),
		reader: reader,
		writer: writer,
	}

	client.Listen()

	return client
}