package protobuf

import (
	pb "../tutorial"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
)

func RunProtoExample() {
	p := &pb.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.Person_HOME},
		},
	}

	out, err := proto.Marshal(p)

	fmt.Println(out, err)

	pp := &pb.Person{}

	if error := proto.Unmarshal(out, pp); error != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	fmt.Println(pp)
}
