import { describe, it, expect } from "vitest";
import "../fixtures";
import { Batch } from "model/batch-edit";

describe("model/batch-edit", () => {
  it("should return defaults", () => {
    const b = new Batch();
    const d = b.getDefaults();
    expect(Array.isArray(d.models)).toBe(true);
    expect(d.values).toEqual({});
    expect(Array.isArray(d.selection)).toBe(true);
  });

  it("should return default form data", () => {
    const b = new Batch();
    const f = b.getDefaultFormData();
    const expectedKeys = [
      "Title",
      "DetailsSubject",
      "Caption",
      "Day",
      "Month",
      "Year",
      "TimeZone",
      "Country",
      "Altitude",
      "Lat",
      "Lng",
      "DetailsArtist",
      "DetailsCopyright",
      "DetailsLicense",
      "DetailsKeywords",
      "Type",
      "Scan",
      "Private",
      "Favorite",
      "Panorama",
      "Iso",
      "FocalLength",
      "FNumber",
      "Exposure",
      "CameraID",
      "LensID",
      "Albums",
      "Labels",
    ];

    expect(Object.keys(f).sort()).toEqual(expectedKeys.sort());
    expect(f.Albums).toEqual({ action: "none", mixed: false, items: [] });
    expect(f.Labels).toEqual({ action: "none", mixed: false, items: [] });
  });

  it("should set selections", () => {
    const b = new Batch();
    b.setSelections([1, 2, 3]);
    expect(b.selection).toEqual([
      { id: 1, selected: true },
      { id: 2, selected: true },
      { id: 3, selected: true },
    ]);
  });

  it("should report selection state for a given id", () => {
    const b = new Batch();
    b.setSelections([1, 2]);
    expect(b.isSelected(1)).toBe(true);
    // toggle one and check again
    b.toggle(1);
    expect(b.isSelected(1)).toBe(false);
    // unknown id returns null per implementation
    expect(b.isSelected(999)).toBeNull();
  });

  it("should toggle and toggleAll", () => {
    const b = new Batch();
    b.setSelections([11, 12, 13]);
    expect(b.getLengthOfAllSelected()).toBe(3);
    b.toggle(12);
    expect(b.isSelected(12)).toBe(false);
    expect(b.getLengthOfAllSelected()).toBe(2);

    b.toggleAll(false);
    expect(b.getLengthOfAllSelected()).toBe(0);

    b.toggleAll(true);
    expect(b.getLengthOfAllSelected()).toBe(3);
  });

  it("should call save and update values from response", async () => {
    const b = new Batch();
    const selection = [5, 7];
    const values = { Title: { value: "New" } };

    // Mock endpoint expected by $api: baseURL is "/api/v1"
    const { Mock } = await import("../fixtures");
    Mock.onPost("api/v1/batch/photos/edit", { photos: selection, values }).reply(
      200,
      { values: { Title: { value: "Saved" } } },
      { "Content-Type": "application/json; charset=utf-8" }
    );

    const result = await b.save(selection, values);
    expect(result).toBe(b);
    expect(b.values).toEqual({ Title: { value: "Saved" } });
  });

  it("should load data (models and values) via getData", async () => {
    const b = new Batch();
    const selection = [101, 102];

    // Response should include models and values
    const { Mock } = await import("../fixtures");
    Mock.onPost("api/v1/batch/photos/edit", { photos: selection }).reply(
      200,
      {
        models: [
          { ID: 1, UID: "ph1", Title: "A" },
          { ID: 2, UID: "ph2", Title: "B" },
        ],
        values: { Title: { mixed: true } },
      },
      { "Content-Type": "application/json; charset=utf-8" }
    );

    await b.getData(selection);

    expect(Array.isArray(b.models)).toBe(true);
    expect(b.models.length).toBe(2);
    expect(b.values).toEqual({ Title: { mixed: true } });
    expect(b.selection).toEqual([
      { id: 101, selected: true },
      { id: 102, selected: true },
    ]);
  });
});
