package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type UserJob struct {
	Role   string
	Squad  string
	Joined time.Time
}

type User struct {
	ID       uint64
	Name     string
	Pronouns []string
	Location string
	Bio      string

	Job UserJob
}

func main() {
	me := map[string]interface{}{
		"id":       uint64(1),
		"name":     "Suhail Patel",
		"pronouns": []string{"he", "him", "his"},
		"location": "London, UK",
		"bio":      "I look at 📉 charts",

		"job": map[string]interface{}{
			"role":   "Backend Engineer",
			"squad":  "Platform",
			"joined": time.Date(2018, 7, 2, 7, 0, 0, 0, time.UTC),
		},
	}

	h := ToUserUsingHandRolled(me)
	fmt.Printf("Hand Rolled User: %+v (User Job: %+v)\n", h, h.Job)

	r, err := ToUserUsingReflect(me)
	if err != nil {
		panic(err) // something bad happened!
	}
	fmt.Printf("Reflect User: %v (User Job: %+v)\n", r, r.Job)
}

// ToUserUsingHandRolled is a function which decodes a map of data using
// manual hand-rolled decoding into the User struct
func ToUserUsingHandRolled(in map[string]interface{}) User {
	u := User{}
	if id, ok := in["id"]; ok {
		u.ID = id.(uint64)
	}
	if name, ok := in["name"]; ok {
		u.Name = name.(string)
	}
	if pronouns, ok := in["pronouns"]; ok {
		u.Pronouns = pronouns.([]string)
	}
	if location, ok := in["location"]; ok {
		u.Location = location.(string)
	}
	if bio, ok := in["bio"]; ok {
		u.Bio = bio.(string)
	}
	if job, ok := in["job"]; ok {
		result := ToUserJobUsingHandRolled(job.(map[string]interface{}))
		u.Job = result
	}
	return u
}

// ToUserJobUsingHandRolled is a function which decodes a map of data
// using manual hand-rolled decoding into the UserJob struct
func ToUserJobUsingHandRolled(in map[string]interface{}) UserJob {
	uj := UserJob{}
	if role, ok := in["role"]; ok {
		uj.Role = role.(string)
	}
	if squad, ok := in["squad"]; ok {
		uj.Squad = squad.(string)
	}
	if joined, ok := in["joined"]; ok {
		uj.Joined = joined.(time.Time)
	}
	return uj
}

// ToUserUsingReflect is a function which uses reflection to dynamically
// assign fields from the input map into struct fields at runtime
func ToUserUsingReflect(in map[string]interface{}) (User, error) {
	u := User{}
	err := UnmarshalUsingReflect(in, &u)
	return u, err
}

// UnmarshalUsingReflect takes in our map and the target struct pointer
// and uses runtime reflection to decode fields into target
func UnmarshalUsingReflect(in map[string]interface{}, target interface{}) error {
	// Make sure the target type is a pointer (to ensure we have passed
	// a reference and not a value)
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("expected a pointer to a struct, got %T", target)
	}

	// Also make sure we have a struct type underpinning the pointer. Note
	// that this struct may not have had any memory allocated for it yet
	targetStruct := reflect.ValueOf(target).Elem()
	if reflect.TypeOf(targetStruct).Kind() != reflect.Struct {
		return fmt.Errorf("expected a pointer to a struct, got %T", target)
	}

	// Final sanity check, make sure our target struct is addressable
	if !targetStruct.CanAddr() {
		return fmt.Errorf("expected an addressable struct for %v", targetStruct)
	}

	structValue := reflect.ValueOf(target).Elem()
	structType := reflect.TypeOf(target).Elem()

	// go through each of our fields in the struct and give
	// it the corresponding lowercase field value from the map
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i) // reflect.StructField
		if val, ok := in[strings.ToLower(field.Name)]; ok {
			fieldValue := structValue.Field(i)

			// If we encounter a map[string]interface{}, then recursively try and
			// unmarshal it into the matching field. The recursive call will handle
			// any pointer allocation we need to deal with
			if m, ok := val.(map[string]interface{}); ok {
				var embeddedTarget interface{}

				switch fieldValue.Kind() {
				case reflect.Ptr:
					if fieldValue.IsNil() {
						// If the struct is nil, we need to allocate that memory
						// beforehand using reflect.New. Note that reflect.New gives
						// a pointer so you want to call it on the element type
						fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
					}
					embeddedTarget = fieldValue.Interface()
				case reflect.Struct:
					// If we have a normal struct, we need to pass down the address
					// (remember, the target struct needs to be addressable)
					embeddedTarget = fieldValue.Addr().Interface()
				default:
					return fmt.Errorf("can only decode inner maps into structs")
				}

				// Recursively unmarshal the struct field
				err := UnmarshalUsingReflect(m, embeddedTarget)
				if err != nil {
					return err
				}
				continue
			}

			fieldValue.Set(reflect.ValueOf(val))
		}
	}

	return nil
}
