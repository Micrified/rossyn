package main

import (
	"fmt"
	"math/rand"
	"math"
	"time"
)


/*
 *******************************************************************************
 *                              Type Definitions                               *
 *******************************************************************************
*/

// Describes a vertex
type Vertex struct {
	chain_id  int                // Identifies the chain this vertex belongs to
	chain_off int                // Offset in the chain (index, effectively) 
	name      string             // Unique name belonging to the callback
	is_sync   bool               // Whether or not this vertex has a sync property
}

// Aliases a list of vetex chains
type Chains [][]*Vertex

// Describes rules for the generation of vertex chains
type Rules struct {
	chain_count       int        // Number of chains to create
	chain_mean_length int        // Mean length (integral) of a chain of vertices
	chain_variance    float64    // Variance in length as % of range (0,chain_mean_length]
	p_callback_merge  float64    // Probability that any two vertices are merged
	p_callback_sync   float64    // Probability that merged vertices must be synchronized
}

// Holds vertices that share common data
type Node struct {
	name      string             // Unique name belonging to node
	vertices  []*Vertex          // List of vertices belonging to the node
}

// Holds nodes that belong to the same executing entity
type Executor struct {
	name      string             // Unique name belonging to executor 
	nodes     []*Node            // List of nodes belonging to the executor
}

// Holds an application
type Application struct {
	name       string            // Name of the application
	executors  []*Executor       // List of executors belonging to application
}

// Enumeration type: Policies (various)
type Policy int 
const (
	Random Policy = iota          // Executor: Assign vertices to random executors
	Complete                      // Executor: Assign vertices across all executors
	Cluster                       // Node: Put vertices from chains in common nodes
	Individual                    // Node: Put vertices in their own nodes
)

// Describes rules for application setup
type Setup struct {
	executor_count    int           // Number of executors
	executor_policy   Policy        // Policy for organizing vertex-chains in executors
	node_policy       Policy        // Policy for organizing vertex-chains in nodes
}

/*
 *******************************************************************************
 *                        Function Definitions: Chains                         *
 *******************************************************************************
*/

// Returns boolean true if two nodes are the same (not a pointer comparison)
func is_same_vertex (a, b *Vertex) bool {
	return (a.chain_id == b.chain_id) && (a.chain_off == b.chain_off)
}

// Outputs an app to stdout
func show_chains (cs *Chains) {
	for i := 0; i < len(*cs); i++ {
		fmt.Printf("%d. ", i)
		for j := 0; j < len((*cs)[i]); j++ {
			if nil == (*cs)[i][j] {
				fmt.Printf("nil-")
			} else {
				fmt.Printf("[%d.%d]-", (*((*cs)[i][j])).chain_id, (*((*cs)[i][j])).chain_off)				
			}
		}
		fmt.Printf("\n")
	}
}

// Generate a application in accordance with the provided rules
func make_chains (r Rules) Chains {
	var chains Chains = make([][]*Vertex, r.chain_count)
	v := float64(r.chain_mean_length) * r.chain_variance
	//fmt.Printf("v = %f\n", v)

	// Seed the PRNG
	rand.Seed(time.Now().UnixNano())

	// Make all chains of length determined by the variance (but always >= 1)
	for i := 0; i < r.chain_count; i++ {
		v_len := -int(math.Round(v)) + 2 * int(math.Round(rand.Float64() * float64(v)))
		//fmt.Printf("v_len(%d) = %d\n", i, v_len)
		c_len := int(math.Max(1.0, float64(r.chain_mean_length + v_len)))
		//fmt.Printf("c_len(%d) = %d\n", i, c_len)
		chains[i] = make([]*Vertex, c_len)

		// Fill chain with vertices
		for j := 0; j < c_len; j++ {
			v_name := fmt.Sprintf("cb_%d_%d", i, j)
			chains[i][j] = &Vertex{chain_id: i, chain_off: j, name: v_name, is_sync: false}
		}
	}

	// Merge (up) possible chain sources using the probability
	for i := 0; i < (r.chain_count - 1); i++ {
		for j := i + 1; j < r.chain_count; j++ {
			if rand.Float64() < r.p_callback_merge {
				fmt.Printf("Source Merge (chain %d into %d)\n", j, i)
				chains[j][0] = chains[i][0]
			}
		}
	}

	// Merge possible chain vertices (but not callbacks)
	for i := 0; i < (r.chain_count - 1); i++ {
		for j := i + 1; j < r.chain_count; j++ {
			for p := 1; p < len(chains[i]); p++ {
				for q := 1; q < len(chains[j]); q++ {
					if rand.Float64() < r.p_callback_merge {

						// No merging if last node merged with same one from p 
						if is_same_vertex(chains[i][p], chains[j][q-1]) {
							continue
						}

						// Possibly a sync node
						if rand.Float64() < r.p_callback_sync {
							chains[i][p].is_sync = true
						}

						// Merge up if (1) no last merge (2) 
						chains[j][q] = chains[i][p]
					}
				}
			}
		}
	}

	// Return app
	return chains
}

/*
 *******************************************************************************
 *                         Function Definitions: Setup                         *
 *******************************************************************************
*/

// Displays an executor on stdout
func show_executors (es *[]*Executor) {
	for _, e := range (*es) {
		fmt.Printf("%s {\n", (*e).name)
		for _, n := range (*e).nodes {
			fmt.Printf("\t%s {", (*n).name)
			for k, v := range (*n).vertices {
				fmt.Printf("[%d.%d]", (*v).chain_id, (*v).chain_off)
				if k < (len((*n).vertices) - 1) {
					fmt.Printf(", ")
				}
			}
			fmt.Println("}")
		}
		fmt.Println("}")
	}
}

// Organizes and generates executors
func make_executors (cs *Chains, setup Setup) []*Executor {
	var executors []*Executor = make([]*Executor, setup.executor_count)

	// Setup executors (create a node for each chain)
	for i := 0; i < setup.executor_count; i++ {
		executor_name := fmt.Sprintf("executor_%d", i)
		executor_nodes := make([]*Node, len((*cs)))

		// Create container nodes for each chain in the executor
		for j := 0; j < len((*cs)); j++ {
			node_name := fmt.Sprintf("node_chain_%d", j)
			executor_nodes[j] = &Node{name: node_name, vertices: []*Vertex{}}
		}

		// Create the executor entity
		executors[i] = &Executor{name: executor_name, nodes: executor_nodes}
	}

	// Distribute vertices from chains across the executors
	// TODO: Implement the actual policies
	// NOTE: Since chains can share vertices, avoid duplicates by 
	// not assigning vertices that don't belong to the chain since
	// they will be assigned by the owner chains
	// [][]*Vertex

	for i, k := 0, 0; i < len(*cs); i++ {
		for j := 0; j < len((*cs)[i]); j++ {

			// If the element in chain i doesn't belong - don't assign
			if (*cs)[i][j].chain_id != i {
				continue
			}

			// Otherwise: Extract executor, and node for chain
			e := executors[k % setup.executor_count]
			n := (*e).nodes[i]

			// Insert vertex into node
			(*n).vertices = append((*n).vertices, (*cs)[i][j])

			// Increment executor that next vertex will be assigned to
			k++
		}
	}

	return executors
}


func main () {
	r := Rules{chain_count: 3, chain_mean_length: 6, chain_variance: 0.5, p_callback_merge: 0.2, p_callback_sync: 0.0}
	cs := make_chains(r)
	fmt.Printf("Done - showing ...\n")
	show_chains(&cs)

	fmt.Printf("Executors ...\n")
	s := Setup{executor_count: 2, executor_policy: Complete, node_policy: Cluster}
	es := make_executors(&cs, s)

	show_executors(&es)
}