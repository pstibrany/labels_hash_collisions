//go:build stringlabels

package main

import "github.com/prometheus/prometheus/model/labels"

var collisions = [][2]labels.Labels{
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl", "HFnEaGl"),
		labels.FromStrings(labels.MetricName, "metric", "lbl", "RqcXatm"),
	},
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl", "gfIS7Ce"),
		labels.FromStrings(labels.MetricName, "metric", "lbl", "x5tSfjf"),
	},
}
