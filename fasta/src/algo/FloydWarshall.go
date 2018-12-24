package algo

import (
    "../structs"
    "../util"
)

/* --- */

// Graph structure for transmitting data between algorithm methods.
type Graph struct {
    segs []structs.Segment /* Segments representing nodes */
    d    [][]int           /* Adjacency matrix and subsequently matrix of greatest path score */
    p    [][]int           /* Matrix of ancestors, i.e. for each pair stores 'middle' node index */
}

/* --- */

// Graph algorithm adapted for program needs. Calculates maximal distances between each pair of nodes.
// Filters and returns segments which form maximal (by score) path.
// Path consists of segments and their connections via gaps.
func FloydWarshall(segs []structs.Segment, gapPenalty, maxDistError int) (int, []structs.Segment ) {
    graph := prepareGraph(segs, gapPenalty, maxDistError)
    floydWarshall(graph)
    return recoverSegs(graph)
}

// Prepares graph by given input segments.
// Each segment forms a node, and edges are gaps connecting segments.
// Note that initially edge's weight is gap penalty + scores of origin segment.
func prepareGraph(segs []structs.Segment, gapPenalty, maxDistErr int) *Graph {

    // Precalculating segments start and end points

    startPoints := make([]structs.Point, len(segs))
    endPoints   := make([]structs.Point, len(segs))

    for i, seg := range segs {
        startPoints[i] = seg.GetStartPoint()
        endPoints[i]   = seg.GetEndPoint()
    }

    //

    d := make([][]int, len(segs))
    p := make([][]int, len(segs))

    for i, seg1 := range segs {
        d[i] = make([]int, len(segs))
        p[i] = make([]int, len(segs))

        for j := range segs {
            p[i][j] = -1

            if i == j {
                d[i][j] = 0
            } else {
                deltaX := util.Abs(endPoints[i].X - startPoints[j].X)
                deltaY := util.Abs(endPoints[i].Y - startPoints[j].Y)

                if startPoints[i].X < startPoints[j].X && startPoints[i].Y < startPoints[j].Y &&
                    deltaX <= maxDistErr && deltaY <= maxDistErr {

                    dist := deltaX + deltaY
                    d[i][j] = seg1.Score + dist * gapPenalty
                } else {
                    d[i][j] = -1000000
                }
            }
        }
    }

    return &Graph { segs, d, p }
}

// Algorithm itself. Is an adopted version of original Floyd-Warshall algorithm.
// For each node pair fills its maximal score and best k-th iteration related to maximal score path.
// Note that graph does not have cycles.
func floydWarshall(graph *Graph) {
    for k := range graph.segs {
        for i := range graph.segs {
            for j := range graph.segs {
                if graph.d[i][j] < graph.d[i][k] + graph.d[k][j] {
                    graph.d[i][j] = graph.d[i][k] + graph.d[k][j]
                    graph.p[i][j] = k
                }
            }
        }
    }
}

// Recovers segments which are in maximal score path.
// Nodes must be filled by algorithm pass first to perform this task.
func recoverSegs(graph *Graph) (int, []structs.Segment) {

    // Search for best score path between origin and destination segments
    maxPair := structs.Point{ 0, 0 }
    maxScore := 0

    for i := range graph.segs {
        for j := range graph.segs {
            // Add missing edge target segment score (should be thought to understand, e.g. on two segments)
            graph.d[i][j] += graph.segs[j].Score

            if graph.d[i][j] > graph.d[maxPair.X][maxPair.Y] {
                maxPair = structs.Point{ i, j }
                maxScore = graph.d[i][j]
            }
        }
    }

    // Fill segments consisted in maximal path

    segsInPath := make([]structs.Segment, 1, 10)
    segsInPath[0] = graph.segs[maxPair.X]

    if maxPair.X != maxPair.Y {
        segsInPath = append(segsInPath, graph.segs[maxPair.Y])
    }

    recoverRec(graph, &segsInPath, maxPair.X, maxPair.Y)

    // Return result

    return maxScore, segsInPath
}

// Segments recover recursive call. Takes k = p[i][j] and splits task into two subcalls (i, j) and (k, j).
func recoverRec(graph *Graph, buffer *[]structs.Segment, i, j int) {
    k := graph.p[i][j]

    if k > -1 {
        *buffer = append(*buffer, graph.segs[k])
        recoverRec(graph, buffer, i, k)
        recoverRec(graph, buffer, k, j)
    }
}
