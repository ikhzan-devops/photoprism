import "../fixtures";
import * as can from "common/can";

let chai = require("chai/chai");
let assert = chai.assert;

describe("common/can", () => {
  it("useVideo", () => {
    assert.equal(can.useVideo, true);
  });

  it("useAVC", () => {
    assert.equal(can.useAVC, true);
  });

  it("useHEVC", () => {
    assert.equal(can.useHEVC, false);
  });

  it("useHEV1", () => {
    assert.equal(can.useHEV1, false);
  });

  it("useVVC", () => {
    assert.equal(can.useVVC, false);
  });

  it("useOGV", () => {
    assert.equal(can.useOGV, true);
  });

  it("useVP8", () => {
    assert.equal(can.useVP8, true);
  });

  it("useVP9", () => {
    assert.equal(can.useVP9, true);
  });

  it("useAV1", () => {
    assert.equal(can.useAV1, true);
  });

  it("useWebM", () => {
    assert.equal(can.useWebM, true);
  });
});
