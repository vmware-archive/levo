/* Copyright (C) 2014 Pivotal Software, Inc.

All rights reserved. This program and the accompanying materials
are made available under the terms of the under the Apache License,
Version 2.0 (the "License”); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.*/
package main

import (
	"code.google.com/p/go.tools/go/vcs"
	"errors"
	"fmt"
	"github.com/cfmobile/levolib"
	"os"
	"os/user"
	"strings"
)

func addTemplatePath(context *levo.Context, templatePath string) ([]levo.TemplateInfo, error) {
	fmt.Printf("")
	var templates []levo.TemplateInfo
	var err error

	fileInfo, err := os.Stat(templatePath)
	if err != nil {
		return []levo.TemplateInfo{}, err
	}

	if fileInfo.IsDir() {
		templates, err = context.AddTemplateDirectory(templatePath)
		if err != nil {
			return []levo.TemplateInfo{}, err
		}
	} else {
		templateObj, err := context.AddTemplateFilePath(templatePath)
		if err != nil {
			return []levo.TemplateInfo{}, err
		}
		templateObj.Directory = "" //User gave us a path to template. Assume no directory information
		templates = make([]levo.TemplateInfo, 0)
		templates = append(templates, templateObj)
	}
	return templates, nil
}

func getUpdatedTemplateRepo(templatePath string) (string, error) {
	if strings.HasPrefix(templatePath, "github.com/") {
		templatePathParts := strings.Split(templatePath, "/")
		rootPrefix, err := getTemplateRepo(strings.Join(templatePathParts[0:3], "/"))
		if err != nil {
			return "", errors.New("Template Repo: " + err.Error())
		}
		templatePath = rootPrefix + templatePath
	}
	return templatePath, nil
}

func getTemplateRepo(repoPath string) (string, error) {
	var (
		vcsCmd         *vcs.Cmd
		repo, rootPath string
	)
	repoRoot, err := vcs.RepoRootForImportPath(repoPath, false)
	if err != nil {
		return "", err
	}
	vcsCmd, repo, rootPath = repoRoot.VCS, repoRoot.Repo, repoRoot.Root
	homeDir, err := userHomeDir()
	if err != nil {
		return "", err
	}
	root := homeDir + "/.levo/" + rootPath
	st, err := os.Stat(root)
	if err == nil && !st.IsDir() {
		return "", errors.New(root + " exists but is not a directory")
	}
	if err != nil {
		err = vcsCmd.Create(root, repo)
		if err != nil {
			return "", err
		}
	} else {
		err = vcsCmd.Download(root)
		if err != nil {
			return "", err
		}
	}

	return homeDir + "/.levo/", nil
}

func userHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
