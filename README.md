# gitb

A command line tool for using Backlog's git comfortably.

## Example

### Show Tree

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/tree/trunk
# default is current
$ gitb tree $TARGET(branch or tag or hash)
```

### Show Commits

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/history/${TARGET}
# default is current
$ gitb commit $TARGET(branch or tag or hash)
```

### Show Branches

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/branches
$ gitb branch
```

### Show Tags

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/tags
$ gitb tag
```

### Show Pull Requests

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/pullRequests?q.statusId=1
$ gitb pr

# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/pullRequests
$ gitb pr -a
```

### Add Pull Request

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/pullRequests/add/...${CURRENT}
$ gitb pr --add
```

### Compare Specific Revision

```
# https://${SPACE_KEY}.backlog.com/git/${PROJECT_KEY}/${REPO_NAME}/compare/${BASE}...${CURRENT}
$ git compare ${BASE} ${CURRENT}
```