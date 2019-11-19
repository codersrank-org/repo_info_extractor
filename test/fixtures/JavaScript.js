// Require - single quotes without semicolon
const a = require('lib1')
require('lib2')
const b = require('lib3')
require('lib4')
const c = require('./lib5')
require('./lib6')

// Require - double quotes without semicolon
const a = require("lib7")
require("lib8")
const b = require("lib9")
require("lib10")
const c = require("./lib11")
require("./lib12")

// Require - single quotes with semicolon
const a = require('lib13');
require('lib14');
const b = require('lib15');
require('lib16');
const c = require('./lib17');
require('./lib18');

// Require - double quotes with semicolon
const a = require("lib19");
require("lib20");
const b = require("lib21");
require("lib22");
const c = require("./lib23");
require("./lib24");

// Import - single quotes without semicolon
import { something } from 'lib25'
import * as whatever from 'lib26'
import whatever from 'lib27'

// Import - double quotes without semicolon
import { something } from "lib28"
import * as whatever from "lib29"
import whatever from "lib30"

// Import - single quotes with semicolon
import { something } from 'lib31';
import * as whatever from 'lib32';
import whatever from 'lib33';

// Import - double quotes with semicolon
import { something } from "lib34";
import * as whatever from "lib35";
import whatever from "lib36";

// Import the entire module
import 'lib37'
import "lib38"
import 'lib39';
import "lib40";
import ('lib41')
import ("lib42")
import ('lib43');
import ("lib44");

import'lib45'
import"lib46"
import'lib47';
import"lib48";
import('lib49')
import("lib50")
import('lib51');
import("lib52");

import     'lib53'
