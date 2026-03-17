#!/bin/bash
# vCluster 创建脚本

TEAMS=("team-a" "team-b")

for team in "${TEAMS[@]}"; do
  echo "Creating vCluster for $team..."
  vcluster create ${team}-cluster \
    --namespace ${team} \
    --create-namespace
done

echo "All vClusters created!"
