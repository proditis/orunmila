# Orunmila
a simple system to refine and produce lists for your bugbounty and pen-test engagements.

The idea behind it is fairly simple a small sqlite(??) database with each word reflecting categories that it can be applied to (ie API, nginx, PHP etc)
This way we can have a huge collection and ask it to give only the given "words" that seem to match your current engagement.

Some imaginary use cases
```sh
getwords -tags nginx,soap,swift,api,xml >wordlist.txt
```

The tool could optionally support extra word databases ie gathered from specific programs.
```
getwords -db programXYZ.db -tags nginx,soap,swift,api,xml
```

The is also a need to be able to import and tag lists easy initially.
```
grabwords -db programXYZ.db -tags generic file_to_extract_words
```

Continued imports of the same file can be utilized to add tags for the words on our database matching the given file (assuming we have it imported already from the previous run under the `generic` tag)
```
grabwords -db programXYZ.db -tags applicationABC file_to_extract_words
```
