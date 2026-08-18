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
