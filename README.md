# Datoms is an implementation of Datomic in Go

## Changes with respect to Datomic

- The zero value is not a valid entity ID. This change has been made to make the zero value more useful. If you transact an entity with entity ID 0 you will create a new entity.

- The zero value is used to communicate omission where possible. For example, there's no need to define entities that use string pointers to signify that the string is optional, you can do that but the empty string value is already treated as if it was unset. The same goes for slices. If a value can be nil it will be treated as optional. This means that if you try to transact the empty string it will not get stored.
