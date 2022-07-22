package goapstra

import (
	"context"
	"log"
	"testing"
)

func TestGetAllTags(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	idList, err := client.listAllTags(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	tagList, err := client.GetAllTags(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(idList) != len(tagList) {
		t.Fatalf("got %d tag IDs but %d tags", len(idList), len(tagList))
	}

	for _, id := range idList {
		tag, err := client.GetTag(context.TODO(), id)
		if err != nil {
			t.Fatal(err)
		}
		log.Println(tag)
	}
}

func TestCreateGetDeleteTag(t *testing.T) {
	client, err := newLiveTestClient()
	if err != nil {
		t.Fatal(err)
	}

	label := TagLabel(randString(10, "hex"))
	description := randString(10, "hex")
	id, err := client.CreateTag(context.TODO(), &DesignTag{
		Label:       label,
		Description: description,
	})

	tag, err := client.GetTag(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}

	if tag.Label != label {
		t.Fatalf("label: '%s' != '%s'", tag.Label, label)
	}

	if tag.Description != description {
		t.Fatalf("description: '%s' != '%s'", tag.Description, description)
	}

	err = client.DeleteTag(context.TODO(), id)
	if err != nil {
		log.Fatal(err)
	}
}
