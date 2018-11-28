package main

import (
	// "fmt"
	"gopkg.in/ldap.v2"
	"gopkg.in/logger.v1"
)

//Group .
type Group struct {
	Cn string
}

func main() {
	l, err := ldap.Dial("tcp", "ldap.xxx.com:389")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	err = l.Bind("cn=xxx,dc=xxx,dc=xxx", "")
	if err != nil {
		panic(err)
	}
	// searchRequest := ldap.NewSearchRequest(
	// 	"ou=users,dc=changhong,dc=com",
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	"(&(objectClass=organizationalPerson)(cn=姓名))",
	// 	[]string{"cn", "uid", "mail", "userPassword"},
	// 	nil,
	// )
	// sr, err := l.Search(searchRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // var matrixs []*types.Matrix
	// for _, entry := range sr.Entries {
	// 	entry.Print()
	// }
	// fmt.Println("len: ", len(sr.Entries))
	dn := "uid=xxx,ou=xxx,dc=xxx,dc=com"
	ar := ldap.NewAddRequest(dn)
	ar.Attribute("objectClass", []string{"person", "top", "organizationalPerson", "inetOrgPerson"})
	ar.Attribute("sn", []string{"描述"})
	ar.Attribute("cn", []string{"姓名"})
	ar.Attribute("mail", []string{"xxx@qq.com"})
	ar.Attribute("userPassword", []string{"123456"})
	err = l.Add(ar)
	if err != nil {
		log.Error(err)
		return
	}
}
