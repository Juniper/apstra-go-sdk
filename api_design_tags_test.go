//go:build integration
// +build integration

package goapstra

import (
	"context"
	"errors"
	"log"
	"testing"
)

func TestGetAllTags(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {

		log.Printf("testing listAllTags() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		idList, err := client.client.listAllTags(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetAllTags() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		tagList, err := client.client.GetAllTags(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		if len(idList) != len(tagList) {
			t.Fatalf("got %d tag IDs but %d tags", len(idList), len(tagList))
		}

		for _, id := range idList {
			log.Printf("testing GetTag() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
			tag, err := client.client.GetTag(context.TODO(), id)
			if err != nil {
				t.Fatal(err)
			}
			log.Println(tag)
		}
	}
}

func TestCreateGetDeleteTag(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for clientName, client := range clients {
		label := TagLabel(randString(10, "hex"))
		description := randString(10, "hex")
		log.Printf("testing CreateTag() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		id, err := client.client.CreateTag(context.TODO(), &DesignTag{
			Label:       label,
			Description: description,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Printf("testing GetTag() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		tag, err := client.client.GetTag(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}

		if tag.Label != label {
			t.Fatalf("label: '%s' != '%s'", tag.Label, label)
		}

		if tag.Description != description {
			t.Fatalf("description: '%s' != '%s'", tag.Description, description)
		}

		log.Printf("testing DeleteTag() against %s %s (%s)", client.clientType, clientName, client.client.ApiVersion())
		err = client.client.DeleteTag(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCreateTagCollision(t *testing.T) {
	clients, err := getTestClients()
	if err != nil {
		t.Fatal(err)
	}

	label := TagLabel(randString(10, "hex"))
	for _, client := range clients {
		id1, err := client.client.CreateTag(context.Background(), &DesignTag{
			Label: label,
		})

		id2, err := client.client.CreateTag(context.Background(), &DesignTag{
			Label: label,
		})
		if err == nil {
			_ = client.client.deleteTag(context.Background(), id1)
			_ = client.client.deleteTag(context.Background(), id2)
			t.Fatal(errors.New("expected error, got none"))
		}
		_ = client.client.deleteTag(context.Background(), id1)

		var ace ApstraClientErr
		if errors.As(err, &ace) && ace.errType == ErrExists { // this is the error we want
			continue
		}
		t.Fatal(err)
	}
}
