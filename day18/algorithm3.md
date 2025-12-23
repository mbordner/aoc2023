




The Shoelace Algorithm (also known as Gauss's area formula or the surveyor's formula) is a method to determine the area of any simple polygon on a grid (coordinate plane) using only the coordinates of its vertices. [1, 2]  
The formula does not require the coordinates to be integers, making it versatile for any polygon where vertices can be described by Cartesian coordinates. [2, 3, 4]  
Steps to Apply the Shoelace Algorithm
To use the algorithm, list the coordinates of the polygon's vertices in order (either clockwise or counter-clockwise). The order determines the sign of the result, which is why the absolute value is taken at the end.

1. List the vertices in a column, repeating the first vertex at the end. For a polygon with vertices $(x_1, y_1), (x_2, y_2), \dots, (x_n, y_n)$: [9, 10]

| X [11] | Y  |
| --- | --- |
| $x_1$ | $y_1$  |
| $x_2$ | $y_2$  |
| ... | ...  |
| $x_n$ | $y_n$  |
| $x_1$ | $y_1$  |

1. Calculate the first sum by multiplying each $x$ coordinate by the $y$ coordinate in the row immediately below it (the "downward" diagonals) and summing the results:$Sum_1 = x_1y_2 + x_2y_3 + \dots + x_ny_1$
2. Calculate the second sum by multiplying each $y$ coordinate by the $x$ coordinate in the row immediately below it (the "upward" diagonals) and summing the results:$Sum_2 = y_1x_2 + y_2x_3 + \dots + y_nx_1$
3. Find the area using the formula:$\text{Area} = \frac{1}{2} |Sum_1 - Sum_2|$ [1, 11, 12]

Example Calculation
For a triangle with vertices (1, 1), (4, 1), and (2, 5): [12]

| X | Y  |
| --- | --- |
| 1 | 1  |
| 4 | 1  |
| 2 | 5  |
| 1 | 1  |

• $Sum_1 = (1 \cdot 1) + (4 \cdot 5) + (2 \cdot 1) = 1 + 20 + 2 = 23$
• $Sum_2 = (1 \cdot 4) + (1 \cdot 2) + (5 \cdot 1) = 4 + 2 + 5 = 11$
• $\text{Area} = \frac{1}{2} |23 - 11| = \frac{1}{2} |12| = 6$ square units [12, 13, 14]

Alternative for Grid Points: Pick's Theorem [15]  
If the polygon's vertices must be on integer lattice points (grid intersections), you can also use Pick's Theorem as a simpler alternative: [16, 17]  
$\text{Area} = I + \frac{B}{2} - 1$

Where:

• $I$ is the number of interior lattice points.
• $B$ is the number of boundary lattice points. [18, 19, 20]

AI responses may include mistakes.

[1] https://www.101computing.net/the-shoelace-algorithm/
[2] https://en.wikipedia.org/wiki/Shoelace_formula
[3] https://11011110.github.io/blog/2021/04/17/picks-shoelaces.html
[4] https://en.wikipedia.org/wiki/Pick%27s_theorem
[5] https://www.sciencedirect.com/science/article/pii/S0020025521012780
[6] https://www.desmos.com/calculator/la7embcze0?lang=it
[7] https://pmc.ncbi.nlm.nih.gov/articles/PMC11398162/
[8] https://www.sciencedirect.com/science/article/pii/S0960148115304390
[9] https://link.springer.com/chapter/10.1007/978-3-031-89873-0_9
[10] https://link.springer.com/chapter/10.1007/978-3-030-88945-6_10
[11] https://courseware.cemc.uwaterloo.ca/42/143/assignments/1140/0
[12] https://www.scribd.com/document/859485377/shoelace
[13] https://www.youtube.com/watch?v=xRfQuk6d5Ic
[14] https://pmc.ncbi.nlm.nih.gov/articles/PMC10975769/
[15] https://www.facebook.com/groups/mathematicsdiscussion/posts/5201735073191105/
[16] https://www.facebook.com/groups/870891947133767/posts/1674862306736723/
[17] https://www.facebook.com/100094746404967/posts/picks-theorempicks-theorem-provides-a-simple-way-to-calculate-the-area-of-a-poly/620382094463332/
[18] https://www3.nd.edu/~ajorza/courses/2018f-m43900/handouts/lecture5.pdf
[19] https://r-knott.surrey.ac.uk/Triangle/TriSolver.html
[20] https://math.hmc.edu/funfacts/picks-theorem/

