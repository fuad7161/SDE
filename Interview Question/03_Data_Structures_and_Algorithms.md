# Data Structures and Algorithms

1. What is an algorithm, and how do you evaluate one?
   - An algorithm is a finite procedure for solving a problem; evaluate its correctness, time and space use, clarity, and behavior on edge cases.
2. What do time complexity and space complexity mean?
   - Time complexity describes how operation count grows with input size, while space complexity describes how memory use grows.
3. What are Big O, Big Omega, and Big Theta notation?
   - Big O gives an asymptotic upper bound, Omega a lower bound, and Theta a tight bound.
4. What are the common complexities from `O(1)` to `O(2^n)`?
   - Typical growth rates are constant, logarithmic, linear, linearithmic, quadratic, polynomial, exponential, and factorial, from generally most to least scalable.
5. What is the difference between an array and a linked list?
   - Arrays provide fast indexed access in contiguous storage; linked lists connect separate nodes and support cheap insertion when the position is already known.
6. When would you use a singly linked list versus a doubly linked list?
   - Use a singly linked list for simpler one-way traversal; use a doubly linked list when backward traversal or easy removal from both directions is needed.
7. What is a stack, and where is it used?
   - A stack is last-in, first-out and is used for call frames, undo operations, parsing, and depth-first search.
8. What is a queue? How do a deque and priority queue differ from it?
   - A queue is first-in, first-out; a deque operates at both ends, while a priority queue removes items according to priority.
9. What is a hash table, and how does it work?
   - It applies a hash function to a key to choose a bucket, giving average `O(1)` lookup, insertion, and deletion with a good distribution.
10. What is a hash collision, and how can it be resolved?
   - A collision occurs when keys map to the same bucket; chaining and open addressing are common resolution techniques.
11. What is the difference between a set and a map?
   - A set stores unique values, while a map associates unique keys with values.
12. What is a tree? Define root, leaf, height, and depth.
   - A tree is a hierarchical acyclic structure; the root has no parent, a leaf no children, depth measures distance from root, and height the longest path to a leaf.
13. What is the difference between a binary tree and a binary search tree?
   - A binary tree allows at most two children; a BST additionally orders keys so searches can discard part of the tree.
14. What is a balanced tree, and why does balance matter?
   - A balanced tree keeps subtree heights close, preserving roughly logarithmic operations instead of degrading to a linear chain.
15. What are preorder, inorder, postorder, and level-order traversal?
   - They visit root-left-right, left-root-right, left-right-root, and breadth by breadth, respectively.
16. What is a heap, and how is it different from a binary search tree?
   - A heap only guarantees that each parent outranks its children and efficiently finds an extreme; a BST maintains ordering for search and traversal.
17. What is a graph? Explain directed, undirected, weighted, and cyclic graphs.
   - A graph contains vertices connected by edges; edges may have direction or weight, and a cyclic graph contains a path returning to its start.
18. What is the difference between breadth-first search and depth-first search?
   - BFS explores by levels using a queue and finds shortest unweighted paths; DFS explores deeply using a stack or recursion.
19. When can binary search be used, and what is its complexity?
   - It requires a sorted, indexable search space and halves that space each step, taking `O(log n)` time.
20. Compare bubble sort, insertion sort, merge sort, and quicksort.
   - Bubble and insertion are usually `O(n²)`; merge sort guarantees `O(n log n)` with extra space, while quicksort averages `O(n log n)` but can reach `O(n²)`.
21. What is a stable sorting algorithm?
   - A stable sort preserves the original relative order of elements with equal keys.
22. What is the difference between recursion, backtracking, and dynamic programming?
   - Recursion is a calling technique, backtracking explores and abandons choices, and dynamic programming saves overlapping subproblem results.
23. What are memoization and tabulation?
   - Memoization is top-down caching of recursive results; tabulation builds results bottom-up in a table.
24. What is a greedy algorithm, and when can it fail?
   - It repeatedly makes the locally best choice; it fails when local choices do not guarantee a global optimum.
25. How would you detect a cycle in a linked list or graph?
   - Use slow/fast pointers for a linked list; for graphs use DFS state or a disjoint set, depending on whether the graph is directed.
26. How would you reverse a linked list?
   - Walk through the nodes while redirecting each `next` pointer to the previous node, keeping temporary references so the remainder is not lost.
27. How would you find duplicates in an array?
   - Track seen values in a set for `O(n)` expected time, or sort first and compare neighbors when extra memory is constrained.
28. How would you find the first non-repeating character in a string?
   - Count each character, then scan the string again and return the first character with count one.
29. How would you check whether two strings are anagrams?
   - Normalize as required, then compare character-frequency maps or sorted characters.
30. How would you find the shortest path in an unweighted graph?
   - Run BFS from the source and store each node's predecessor to reconstruct the first-discovered shortest path.
