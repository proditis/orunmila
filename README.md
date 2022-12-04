# Orunmila
a simple system to refine and produce lists for your bugbounty and pen-test engagements.

The idea behind it is fairly simple a small sqlite(??) database with each word associated tags. Each word in the dictionary can be associated with multiple tags.
This provides for a way to later request the words from a database based on specific tags and use the generated wordlist with you normal tools, be it ffuf, dirbuster etc.


## Building
NOTE: This is still a really early prototype so not much of a build system into the mix.

```sh
go build orunmila
```


## Examples
* Import words from `lista.txt` and tag as `lista`
  ```
  orunmila -tags lista lista.txt
  ```

* List words with tag as `lista`
  ```
  orunmila -tags lista
  ```

* Import words from `listb.txt` and tag as `listb`
  ```
  orunmila -tags listb listb.txt
  ```

* List words with tag as `listb`
  ```
  orunmila -tags listb
  ```

* Import words from `lista.txt` and `listb.txt` and tag as `listc`
  ```
  orunmila -tags listc lista.txt listb.txt
  ```

* List words with tag as `listc` (should return all words)
  ```
  orunmila -tags listc
  ```

### Drupal example
Take the following hypothetical scenario, we have a target system that is based on drupal. We have already populated our `orunmila.db` with appropriate words and tags before hand.

Using orunmila we extract the keywords that match our criteria
```sh
orunmila -tags drupal,dir,nginx,php >drupal_words.txt
ffuf -w drupal_words.txt -u https://drupal-target/FUZZ
```

The tool supports using specific database files ie
```
orunmila -db programXYZ.db -tags nginx,soap,swift,api,xml
```

You can use Orunmila to import wordlists into your database with given set of tags. Existing words will have their tags updated to reflect the new ones
```
orunmila -db programXYZ.db -tags raft,directories,manual raft-medium-directories.txt
```
