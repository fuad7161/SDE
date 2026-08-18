# Git and Version Control

1. What is version control, and why is it useful?
   - Version control records file history so teams can collaborate, review changes, create branches, and restore earlier versions.
2. What is Git, and how is it different from GitHub or GitLab?
   - Git is a distributed version-control tool; GitHub and GitLab host Git repositories and add collaboration and automation features.
3. What are a working tree, staging area, and repository?
   - The working tree contains current files, the staging area selects the next snapshot, and the repository stores committed history.
4. What is a commit?
   - A commit is an immutable snapshot of staged changes with metadata and links to its parent commit or commits.
5. What is a branch, and why do teams use branches?
   - A branch is a movable reference to a commit, allowing work to develop independently before integration.
6. What is `HEAD` in Git?
   - `HEAD` identifies the currently checked-out commit, usually indirectly through the current branch.
7. What is the difference between `git fetch` and `git pull`?
   - `fetch` downloads remote history without integrating it; `pull` fetches and then merges or rebases into the current branch.
8. What is the difference between merging and rebasing?
   - Merging combines histories with a merge commit, while rebasing replays commits on a new base and rewrites their identities.
9. What is a fast-forward merge?
   - It moves the target branch pointer forward when no divergent commits require a merge commit.
10. What is a merge conflict, and how do you resolve one?
   - A conflict means Git cannot safely combine overlapping changes; edit the marked files, stage the resolved versions, and finish the operation.
11. What is the difference between `git clone` and `git fork`?
   - Clone copies a repository locally; a fork is a separate server-side repository derived from another project.
12. What is the purpose of a remote such as `origin`?
   - A remote is a named reference to another repository used to fetch and push history; `origin` is the conventional default name.
13. What is the purpose of `.gitignore`?
   - It specifies intentionally untracked paths Git should normally ignore, such as build output and local configuration.
14. What is the difference between tracked, untracked, and ignored files?
   - Tracked files are in Git history or staging, untracked files are new to Git, and ignored files match ignore rules.
15. What is the difference between `git reset` and `git revert`?
   - Reset moves a branch and may alter local state; revert safely adds a new commit that undoes an earlier commit.
16. What is `git stash`, and when would you use it?
   - Stash temporarily saves uncommitted changes so the working tree can be cleaned and the changes reapplied later.
17. What does `git cherry-pick` do?
   - It applies the changes introduced by selected commits onto the current branch as new commits.
18. What are Git tags used for?
   - Tags give stable names to particular commits, commonly marking releases such as `v2.0.0`.
19. How can you inspect the history of a file or line?
   - Use `git log -- <file>` for file history and `git blame <file>` to identify the last commit affecting each line.
20. What makes a good commit message?
   - It has a concise imperative summary and, when useful, explains the motivation and consequences rather than restating the diff.
21. What is a pull request or merge request?
   - It is a hosted review workflow for discussing, checking, and approving a proposed branch integration.
22. Why should generated files and secrets usually not be committed?
   - Generated files create noise and reproducibility problems; committed secrets can be copied from history even after later deletion.

## Medium to Advanced

23. How does Git represent commits, trees, blobs, and tags internally?
   - **Key note:** Git stores content-addressed objects: blobs hold content, trees directories, commits snapshots/history, and annotated tags references plus metadata.
24. Why does changing a commit create a new commit hash?
   - **Key note:** A commit ID hashes its contents and parent references, so any metadata, tree, or ancestry change creates a new object.
25. What is the difference between a merge commit, squash merge, and rebase merge?
   - **Key note:** They preserve branch topology, combine changes into one commit, or replay individual commits into linear history.
26. When is rebasing a shared branch dangerous?
   - **Key note:** Rebase rewrites commit identities, forcing collaborators to reconcile history they already based work upon.
27. How does `git rebase --onto` help move a range of commits?
   - **Key note:** It replays commits after a selected old base onto a different new base, enabling precise history surgery.
28. What is an interactive rebase used for?
   - **Key note:** It can reorder, edit, squash, split, or remove unpublished commits before sharing them.
29. What is the reflog, and when can it recover lost work?
   - **Key note:** Reflog records local reference movements, allowing recently orphaned commits to be located before expiration.
30. What is a detached `HEAD`, and how do you preserve commits created there?
   - **Key note:** `HEAD` points directly to a commit; create a branch or tag before those commits become unreachable.
31. How do the three forms of `git reset` differ?
   - **Key note:** Soft moves only the branch, mixed also resets staging, and hard also overwrites tracked working-tree content.
32. How do `ours` and `theirs` change meaning during rebase?
   - **Key note:** During rebase, “ours” is the target/upstream side and “theirs” is the commit being replayed.
33. What is a three-way merge?
   - **Key note:** Git compares two tips with their common ancestor to determine each side's independent changes.
34. How does Git detect renames if it does not store rename operations?
   - **Key note:** It compares deleted and added content heuristically based on similarity during diff or merge.
35. What is `git bisect`, and how does it locate a regression?
   - **Key note:** Binary search repeatedly tests midpoint commits between known good and bad revisions.
36. What is the difference between submodules and subtrees?
   - **Key note:** Submodules reference an external repository commit; subtrees copy and merge another project's content into the repository.
37. What problem does Git LFS solve?
   - **Key note:** It stores small pointers in Git while large binary content lives in separate object storage.
38. How can a secret be removed from the entire Git history?
   - **Key note:** Rotate it first, rewrite all affected history with a filtering tool, force-update refs, and coordinate every clone.
39. What are signed commits and signed tags?
   - **Key note:** Cryptographic signatures let others verify the claimed author/key and that tagged or committed data was not changed.
40. How do shallow and partial clones differ?
   - **Key note:** Shallow clones limit history depth; partial clones omit selected object content and fetch it on demand.
41. What is a merge strategy versus a merge strategy option?
   - **Key note:** A strategy chooses the merge algorithm; an option tunes that algorithm's conflict resolution behavior.
42. How would you design a branch-protection policy for a production repository?
   - **Key note:** Require reviewed pull requests, passing checks, protected history, limited bypass, signed artifacts, and controlled release permissions.
