// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     public/app/plugins/gen.go
// Using jennies:
//     TSTypesJenny
//     LatestMajorsOrXJenny
//     PluginEachMajorJenny
//
// Run 'make gen-cue' from repository root to regenerate.

import * as common from '@grafana/schema';

export const pluginVersion = "10.2.0-pre";

export interface Options extends common.SingleStatBaseOptions {
  maxVizHeight: number;
  minVizHeight: number;
  minVizWidth: number;
  showThresholdLabels: boolean;
  showThresholdMarkers: boolean;
}

export const defaultOptions: Partial<Options> = {
  maxVizHeight: 1000,
  minVizHeight: 75,
  minVizWidth: 75,
  showThresholdLabels: false,
  showThresholdMarkers: true,
};
