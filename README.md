# ROSSYN 

ROSSYN, or the ROS Synthesizer, generates applications for ROS2. It has a particular focus on generating applications with variable callback-chains, among other properties. 

## Input(s)

1. **Chain generation rules**
    a. `chain_count`: How many independent chains exist in the application
    b. `chain_mean_length`: How many callbacks chains have on average
    c. `chain_variance`: Given as a percentage of the mean chain length. Controls deviation from mean length up/to +/- 100% of the given mean length
    d. `p_callback_merge`: The probability that any two callbacks in a chain are shared. 
    e. `p_callback_sync`: The probability that a merging callback between two chains require synchronization
2. **Setup generation rules** 
    a. `executor_count`: How many executors should multiplex the callbacks
    b. `executor_policy`: Policy governing how callback chains are distributed across executors
    c. `node_policy`: Policy governing how callbacks from a chain are organized in nodes

## Output(s)

At this time only 


## Notes

Miscellaneous design considerations include the following: 

* Callbacks communicate over topics. These have been omitted from the generation tool thus far, as they can be simply created last minute. However, there remains considerations to make around how topics are chosen with respect to merged and sync nodes. Nodes that are shared between chains may share the same topic. This makes sense in most applications, and so it is assumed (for future work on this tool) that a single topic will be created for all callbacks routed through a merged node. However, sync nodes are special cases where all incident edges must communicate something before the node can propagate output. On a topic level, this means there must exist independent topics over which all incident chains/edges communicate. Thus, sync nodes will have dedicated topics for each incident edge from a chain. 

## Schedule synthesis

This tool is intended to be expanded later with considerations for schedule synthesis. The README will be updated accordingly