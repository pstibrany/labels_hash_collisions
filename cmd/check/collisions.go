//go:build !stringlabels

package main

import "github.com/prometheus/prometheus/model/labels"

var collisions = [][2]labels.Labels{
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl", "qeYKm3"),
		labels.FromStrings(labels.MetricName, "metric", "lbl", "2fUczT"),
	},
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "value", "lbl2", "l6CQ5y"),
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "value", "lbl2", "v7uDlF"),
	},
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "W7qx", "lbl2", "zqqr"),
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "Z00w", "lbl2", "wuwb"),
	},
	{
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "59zo", "lbl2", "ucIY"),
		labels.FromStrings(labels.MetricName, "metric", "lbl1", "ThBT", "lbl2", "XYrv"),
	},
	{
		labels.FromStrings(labels.MetricName, "pqrw", "lbl", "Aanhoh"),
		labels.FromStrings(labels.MetricName, "sBmm", "lbl", "pwdthe"),
	},
}
