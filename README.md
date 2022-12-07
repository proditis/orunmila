# Orunmila
A simple tool to refine and produce lists for your bugbounty and pen-test engagements.

The idea behind it is fairly simple, a small sqlite(??) database with each word associated tags. Each word in the dictionary can be associated with multiple tags.
This provides for a way to later request the words from a database based on specific tags and use the generated wordlist with you normal tools, be it ffuf, dirbuster etc.


## Installation

```sh
GO111MODULE=on go install github.com/proditis/orunmila@latest
```

## Building
NOTE: This is still a really early prototype so not much of a build system into the mix.

```sh
export CGO_ENABLED=1
go get github.com/mattn/go-sqlite3
go get github.com/sirupsen/logrus
go build orunmila.go
```


## Subcommands
* **`add`** words from the cli
  ```sh
  orunmila add -tags a,b,c word1 word2 word3
  ```
* **`import`** words from a file
  ```sh
  orunmila import -tags a,b,c filename
  ```
* **`search`** words
  ```
    orunmila search -tags a,b,c filename
  ```
* **`vacuum`** database and apply any schema updates
  ```sh
  orunmila vacuum a
  ```
* **`describe`** a database
  ```sh
  orunmila describe My Description for this database
  ```
* **`info`** return information about a database
  ```sh
  $ orunmila info all
  [version]: 0.0.0
  [dbname]: default
  [description]: My Description for this database
  ```

## Examples
* Import words from `lista.txt` and tag as `lista`
  ```
  orunmila import -tags lista lista.txt
  ```

* List words with tag as `lista`
  ```
  orunmila search -tags lista
  ```

* Import words from `listb.txt` and tag as `listb`
  ```
  orunmila import -tags listb listb.txt
  ```

* List words with tag as `listb`
  ```
  orunmila search -tags listb
  ```

* Import words from `lista.txt` and `listb.txt` and tag as `listc`
  ```
  orunmila import -tags listc lista.txt listb.txt
  ```

* List words with tag as `listc` (should return all words)
  ```
  orunmila search -tags listc
  ```

### Drupal example
Take the following hypothetical scenario, we have a target system that is based on drupal. We have already populated our `orunmila.db` with appropriate words and tags before hand.

Using orunmila we extract the keywords that match our criteria
```sh
orunmila search -tags drupal,dir,nginx,php >drupal_words.txt
ffuf -w drupal_words.txt -u https://drupal-target/FUZZ
```

The tool supports using specific database files ie
```
orunmila search -db programXYZ.db -tags nginx,soap,swift,api,xml
```

You can use Orunmila to import wordlists into your database with given set of tags. Existing words will have their tags updated to include old and new ones
```
orunmila import -db programXYZ.db -tags raft,directories,manual raft-medium-directories.txt
```

Add a new drupal entry you discovered from the command line (without file)
```
orunmila add -tags drupal,directory word1 word2
```
