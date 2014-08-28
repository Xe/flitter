# Image builder overview

```
git push -> git repo -> builder -> docker image -> deploy to deis
```

## Settings

- [ ] Private registry host/port

## Builder Process

- [ ] Create temporary directory
- [ ] Extract a copy of the branch you want to deploy
- [ ] Grab runtime config from controller
  - [ ] store config in memory
- [ ] Find the Dockerfile
- [ ] Build Docker image
  - [ ] Check for multiple exposed ports, error out if so
  - [ ] splay in config
  - [ ] Tag for private registry
  - [ ] Push to private registry
- [ ] Post information about the build to the controller
  - [ ] collect data into a struct
- [ ] print end message
  - [ ] Show domain name of app
  - [ ] if no domain name
    - [ ] error and show message
- [ ] perform git gc on repo
