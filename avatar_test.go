package main

import (
    "testing"
    "path"
    "os"
    "io/ioutil"
    // "github.com/stretchr/gomniauth/test"
    // "github.com/stretchr/testify"
    gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
    var authAvatar AuthAvatar
    testUser := &gomniauthtest.TestUser{}
    testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
    testChatUser := &chatUser{User: testUser}
    url, err := authAvatar.GetAvatarURL(testChatUser)
    if err != ErrNoAvatarURL {
        t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
    }
    testUrl := "http://url-to-gravatar/"
    testUser = &gomniauthtest.TestUser()
    testChatUser.User = testUser
    testUrl.On("AvatarURL").Return(testChatUser)
    if err != nil {
        t.Error("AuthAvatar.GetAvatarURL should return no error when value is present")
    } else {
        if url != testUrl {
            t.Error("AuthAvatar.GetAvatarURL should return correct URL")
        }
    }
}

/*
func TestAuthAvatar(t *testing.T) {
    var authAvatar AuthAvatar
    client := new(client)
    url, err := authAvatar.GetAvatarURL(client)
    if err != ErrNoAvatarURL {
        t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value is present")
    }
    // set a value
    testUrl := "http://url-to-gravatar/"
    client.userData = map[string]interface{}{"avatar_url": testUrl}
    url, err = authAvatar.GetAvatarURL(client)
    if err != nil {
        t.Error("AuthAvatar.GetAvatarURL should return no error when value is present")
    } else {
        if url != testUrl {
            t.Error("AuthAvatar.GetAvatarURL should return correct URL")
        }
    }
}*/

func TestGravatarAvatar(t *testing.T) {
    var gravatarAvitar GravatarAvatar
    user := &chatUser{uniqueID: "abc"}
    url, err := gravatarAvitar.GetAvatarURL(user)
    if err != nil {
        t.Error("GravatarAvitar.GetAvatarURL should not return an error")
    }
    if url != "//www.gravatar.com/avatar/abc" {
        t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
    }
}

/*
func TestGravatarAvatar(t *testing.T) {
    var gravatarAvitar GravatarAvatar
    client := new(client)
    client.userData = map[string]interface{}{"userid": 
        "0bc83cb571cd1c50ba6f3e8a78ef1346"}
    url, err := gravatarAvitar.GetAvatarURL(client)
    if err != nil {
        t.Error("GravatarAvitar.GetAvatarURL should not return an error")
    }
    if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
        t.Errorf("GravatarAvitar.GetAvatarURL wrongly returned %s", url)
    }
}*/

func TestFileSystemAvatar(t *testing.T) {
    // make a test avatar file
    filename := path.Join("avatars", "abc.jpg")
    ioutil.WriteFile(filename, []byte{}, 0777)
    defer func() { os.Remove(filename) }()
    
    var fileSystemAvatar FileSystemAvatar
    user := &chatUser{uniqueID: "abc"}
    url, err := fileSystemAvatar.GetAvatarURL(user)
    if err != nil {
        t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
    }
    if url != "/avatars/abc.jpg" {
        t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
    }
}

/*
func TestFileSystemAvatar(t *testing.T) {
    // make a test avatar file
    filename := path.Join("avatars", "abc.jpg")
    ioutil.WriteFile(filename, []byte{}, 0777)
    defer func() { os.Remove(filename) }()
    
    var fileSystemAvatar FileSystemAvatar
    client := new(client)
    client.userData = map[string]interface{}{"userid": "abc"}
    url, err := fileSystemAvatar.GetAvatarURL(client)
    if err != nil {
        t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
    }
    if url != "/avatars/abc.jpg" {
        t.Errorf("FileSystemAvatar.GetAvatarURL wrongly returned %s", url)
    }
}*/
