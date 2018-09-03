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
	l, err := ldap.Dial("tcp", "ldap.changhong.com:389")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=changhong,dc=com", "LivVM88va238lmvLI")
	if err != nil {
		panic(err)
	}
	// searchRequest := ldap.NewSearchRequest(
	// 	"ou=users,dc=changhong,dc=com",
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	"(&(objectClass=organizationalPerson)(cn=李蕊材))",
	// 	[]string{"cn", "uid", "mail", "userPassword"},
	// 	//	[]string{"admin", "mce-dev-admin", "mce-prod-admin", "mce-test-admin"},
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

	dn := "uid=mstar.michael,ou=users,dc=changhong,dc=com"

	ar := ldap.NewAddRequest(dn)
	ar.Attribute("objectClass", []string{"person", "top", "organizationalPerson", "inetOrgPerson"})
	ar.Attribute("sn", []string{"Mstar技术支持"})
	ar.Attribute("cn", []string{"肖海军"})
	ar.Attribute("mail", []string{"michael.shaw@mstarsemi.com"})
	ar.Attribute("userPassword", []string{"{SSHA}shxynyMxqRxy1ng/FvIC8lo7g3Qfmro3sUIqVw=="})
	err = l.Add(ar)
	if err != nil {
		log.Error(err)
		return
	}

}
