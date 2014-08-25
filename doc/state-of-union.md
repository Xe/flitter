## How Deis' Builder Works

### On Push

1. Gitrecieve is called with the `run` argument. It parses out the original SSH 
   command and scrapes out the repository being pushed to.
    - If the repository does not exist it creates it on disk and adds pre and 
    post-recieve hooks to call the builder
2. Gitrecieve calls the original SSH command inside `git-shell` and has it 
   process the push, then calling the slugbuilder.
3. The hook calls the receiver check and checks against Deis for authentication 
   to push to the repo.
    - If the user does not have permission they are denied access.
    - If the user does have permission the commits are accepted and added to 
    the server-side repo, then the builder is run.
4. See builder step by step teardown.
