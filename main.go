package main

import (
	"context"
	"bufio"
    "fmt"
    "io"
    "io/ioutil"
    "os"
	"sort"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

func stringToVersionSlice(stringSlice []string) []*semver.Version {
	versionSlice := make([]*semver.Version, len(stringSlice))
	for i, versionString := range stringSlice {
		versionSlice[i] = semver.New(versionString)
	}
	return versionSlice
}

func versionToStringSlice(versionSlice []*semver.Version) []string {
	stringSlice := make([]string, len(versionSlice))
	for i, version := range versionSlice {
		stringSlice[i] = version.String()
	}
	return stringSlice
}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	//first sort the releases, ascending order
	sort.Sort(releases)
	//revert the matrix
	for i, j := 0, len(releases)-1; i < j; i, j = i+1, j-1 {
        releases[i], releases[j] = releases[j], releases[i]
    }

	i = 0
	for releases[i] > minVersion {
		i = i+1
	}

	// return releases
	versionSlice := releases[i:len(releases)]
	return versionSlice
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	//read the user input
	reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter filename: ")
    filename, _ := reader.ReadString('\n')
    //read the file
    fileString, err := ioutil.ReadFile("file.txt") 

    repoUserNames := []string{}
    repoNames := []string{}
    minVersions := []string{}
    i := 0
    temp := ""
    //loop through the file string, storing the data correctly
    for i=0; i<len(fileString); i++{
    	if fileString[i] == ','{
    		repoNames = append(repoNames, temp)
    		temp = ""
    	}else if fileString[i] == '/'{
    		repoUserNames = append(repoNames, temp)
    		temp = ""
    	}else if fileString[i] =='\n'{
    		minVersions = append(minVersions, temp)
    		temp = ""
    	}else{
    		temp = temp + fileString[i]
    	}
    }

    //if the first line is actually repository,min_version, then delete the first value of reponames and minverisons
    repoNames = repoNames[1:len(repoNames)]
    minVersions = minVersions[1:len(minVersions)]

    //convert minversions to semver
    minVersionsSemVer := stringToVersionSlice(minVersions)

	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}

	//now load each repository
	for i=0; i<len(repoUserNames); i++{
		releases, _, err := client.Repositories.ListReleases(ctx, repoUserNames[i], repoNames[i], opt)
		for err != nil { //this will try to load the thing forever, so if the user makes a mistaken input, tough luck
			releases, _, err = client.Repositories.ListReleases(ctx, repoUserNames[i], repoNames[i], opt)
		}

		//get the correct minver
		minVersion := semver.New(minVersionsSemVer[i])

		//change all releases
		allReleases := make([]*semver.Version, len(releases))
		semVerOut := LatestVersions(allReleases,minVersion)

		//output the formatted string
		print("latest versions of %s/%s: %s",repoUserNames[i],repoNames[i], versionToStringSlice(semVerOut))
	}
	//now loop through each version of the file
}
