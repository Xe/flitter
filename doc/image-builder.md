# Image builder overview

```
git push -> git repo -> builder -> docker image -> deploy to deis
```

## Builder Process

- [ ] Create temporary directory
- [ ] Extract a copy of the branch you want to deploy
- [ ] Grab runtime config from controller
  - [ ] store config in memory
- [ ] Find the Docker/Procfile
  - [ ] If Dockerfile:
    - [ ] Fallthru, have Docker build repo raw
  - [ ] If Procfile:
    - [ ] Run Slugbuilder
      - [ ] Save buildpack cache
      - [ ] Mount buidpacks
      - [ ] Mount app to /tmp/app
    - [ ] Attach client to output of Slugbuilder
    - [ ] Extract slug from container
- [ ] Build Docker image
  - [ ] Check for multiple exposed ports, error out if so
  - [ ] splay in config
  - [ ] Tag for private registry
  - [ ] Push to private registry
- [ ] Extract process types from Procfile
  - [ ] store in memory
- [ ] Post information about the build to the controller
  - [ ] collect data into a struct
  - [ ] jsonify it and send to controller
- [ ] print end message
  - [ ] Show domain name of app
  - [ ] if no domain name
    - [ ] error and show message
- [ ] perform git gc on repo

