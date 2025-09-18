import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { shallowMount } from "@vue/test-utils";
import { nextTick } from "vue";
import PPhotoBatchEdit from "component/photo/edit/batch.vue";
import { Batch } from "model/batch-edit";

// Mock the models and dependencies
vi.mock("model/batch-edit");
vi.mock("model/album");
vi.mock("model/label");
vi.mock("model/thumb");

describe("component/photo/edit/batch", () => {
  let wrapper;
  let mockBatch;

  const mockSelection = ["uid1", "uid2", "uid3"];

  const mockModels = [
    {
      ID: 1,
      UID: "uid1",
      Title: "Photo 1",
      FileName: "photo1.jpg",
      Type: "image",
      getOriginalName: () => "photo1.jpg",
      thumbnailUrl: (size) => `/thumb/${size}/photo1.jpg`,
    },
    {
      ID: 2,
      UID: "uid2",
      Title: "Photo 2",
      FileName: "photo2.jpg",
      Type: "video",
      getOriginalName: () => "photo2.jpg",
      thumbnailUrl: (size) => `/thumb/${size}/photo2.jpg`,
    },
    {
      ID: 3,
      UID: "uid3",
      Title: "Photo 3",
      FileName: "photo3.jpg",
      Type: "live",
      getOriginalName: () => "photo3.jpg",
      thumbnailUrl: (size) => `/thumb/${size}/photo3.jpg`,
    },
  ];

  const mockValues = {
    Title: { value: "Test Title", mixed: false },
    Caption: { value: "", mixed: true },
    DetailsSubject: { value: "Test Subject", mixed: false },
    Day: { value: 15, mixed: false },
    Month: { value: 6, mixed: false },
    Year: { value: 2023, mixed: false },
    TimeZone: { value: "UTC", mixed: false },
    Country: { value: "US", mixed: false },
    Altitude: { value: 100, mixed: false },
    Lat: { value: 37.7749, mixed: false },
    Lng: { value: -122.4194, mixed: false },
    DetailsArtist: { value: "Test Artist", mixed: false },
    DetailsCopyright: { value: "Test Copyright", mixed: false },
    DetailsLicense: { value: "Test License", mixed: false },
    Type: { value: "image", mixed: false },
    Scan: { value: true, mixed: false },
    Favorite: { value: false, mixed: true },
    Private: { value: false, mixed: false },
    Panorama: { value: false, mixed: false },
    Albums: { items: [], mixed: false, action: "none" },
    Labels: { items: [], mixed: false, action: "none" },
  };

  beforeEach(() => {
    // Reset mocks
    vi.clearAllMocks();

    // Create mock batch instance
    mockBatch = {
      models: mockModels,
      values: mockValues,
      selection: [
        { id: "uid1", selected: true },
        { id: "uid2", selected: true },
        { id: "uid3", selected: true },
      ],
      getData: vi.fn().mockResolvedValue(mockBatch),
      save: vi.fn().mockResolvedValue(mockBatch),
      getValuesForSelection: vi.fn().mockResolvedValue(mockValues),
      getDefaultFormData: vi.fn().mockReturnValue({
        Title: { value: "Test", action: "none", mixed: false },
        DetailsSubject: { value: "", action: "none", mixed: false },
        Caption: { value: "", action: "none", mixed: false },
        Day: { value: 0, action: "none", mixed: false },
        Month: { value: 0, action: "none", mixed: false },
        Year: { value: 0, action: "none", mixed: false },
        TimeZone: { value: "UTC", action: "none", mixed: false },
        Country: { value: "US", action: "none", mixed: false },
        Altitude: { value: 0, action: "none", mixed: false },
        Lat: { value: 37.7749, action: "none", mixed: false },
        Lng: { value: -122.4194, action: "none", mixed: false },
        DetailsArtist: { value: "", action: "none", mixed: false },
        DetailsCopyright: { value: "", action: "none", mixed: false },
        DetailsLicense: { value: "", action: "none", mixed: false },
        DetailsKeywords: { value: "", action: "none", mixed: false },
        Type: { value: "image", action: "none", mixed: false },
        Iso: { value: 0, action: "none", mixed: false },
        FocalLength: { value: 0, action: "none", mixed: false },
        FNumber: { value: 0, action: "none", mixed: false },
        Exposure: { value: "", action: "none", mixed: false },
        CameraID: { value: 0, action: "none", mixed: false },
        LensID: { value: 0, action: "none", mixed: false },
        Scan: { value: false, action: "none", mixed: false },
        Private: { value: false, action: "none", mixed: false },
        Favorite: { value: false, action: "none", mixed: false },
        Panorama: { value: false, action: "none", mixed: false },
        Albums: { items: [], mixed: false, action: "none" },
        Labels: { items: [], mixed: false, action: "none" },
      }),
      getLengthOfAllSelected: vi.fn().mockReturnValue(3),
      isSelected: vi.fn().mockReturnValue(true),
      toggle: vi.fn(),
      toggleAll: vi.fn(),
    };

    // Mock Batch constructor
    Batch.mockImplementation(() => mockBatch);

    // Mock global objects
    global.$notify = {
      success: vi.fn(),
      error: vi.fn(),
    };

    global.$lightbox = {
      openModels: vi.fn(),
    };

    global.$event = {
      subscribe: vi.fn(),
      unsubscribe: vi.fn(),
    };

    wrapper = shallowMount(PPhotoBatchEdit, {
      props: {
        visible: false, // Start with false to avoid initial rendering issues
        selection: mockSelection,
        openDate: vi.fn(),
        openLocation: vi.fn(),
        editPhoto: vi.fn(),
      },
      global: {
        mocks: {
          $notify: global.$notify,
          $lightbox: global.$lightbox,
          $event: global.$event,
          $config: {
            feature: vi.fn().mockReturnValue(true),
          },
          $vuetify: { display: { mdAndDown: false } },
        },
        stubs: {
          VDialog: {
            template: '<div class="v-dialog">' + '<slot v-if="modelValue" />' + "</div>",
            props: ["modelValue"],
          },
          VDataTable: {
            template: '<div class="v-data-table"></div>',
            props: ["headers", "items"],
          },
          PLocationInput: {
            template: '<div class="p-location-input"></div>',
            props: ["latlng", "label"],
            emits: ["update:latlng", "changed", "open-map", "delete", "undo"],
          },
          PLocationDialog: {
            template: '<div class="p-location-dialog"></div>',
            props: ["visible", "latlng"],
            emits: ["close", "confirm"],
          },
          BatchChipSelector: {
            template: '<div class="batch-chip-selector"></div>',
            props: ["items", "availableItems"],
            emits: ["update:items"],
          },
          IconLivePhoto: {
            template: '<i class="icon-live-photo"></i>',
          },
        },
      },
    });

    // Ensure internal data initialized similar to when visible watcher runs
    // Avoid direct assignment before component is ready; simulate visible=true flow
    wrapper.vm.values = { ...mockValues };
    if (typeof wrapper.vm.setFormData === "function") {
      wrapper.vm.setFormData();
    }
    wrapper.vm.allSelectedLength = mockBatch.getLengthOfAllSelected();
  });

  afterEach(() => {
    if (wrapper) {
      wrapper.unmount();
    }
  });

  describe("Computed Properties", () => {
    beforeEach(() => {
      // Set up component state for computed property tests
      wrapper.vm.model = mockBatch;
      wrapper.vm.values = mockValues;
      // Merge into existing complete formData to avoid template access errors
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Lat: { value: 37.7749, action: "none", mixed: false },
        Lng: { value: -122.4194, action: "none", mixed: false },
      };
    });

    it("should compute form title correctly", () => {1
      expect(wrapper.vm.formTitle).toBe("Batch Edit (3)");
    });

    it("should compute current coordinates correctly", () => {
      const coords = wrapper.vm.currentCoordinates;
      expect(coords).toEqual([37.7749, -122.4194]);
    });

    it("should handle mixed location state", () => {
      wrapper.vm.values = {
        Lat: { mixed: true },
        Lng: { mixed: true },
      };

      expect(wrapper.vm.isLocationMixed).toBe(true);
      expect(wrapper.vm.currentCoordinates).toEqual([0, 0]);
    });
  });

  describe("Form Data Management", () => {
    beforeEach(() => {
      wrapper.vm.model = mockBatch;
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Title: { value: "Changed", action: "update", mixed: false },
        Caption: { value: "Original", action: "none", mixed: false },
      };
    });

    it("should correctly detect unsaved changes true/false", async () => {
      expect(wrapper.vm.hasUnsavedChanges()).toBe(true);
      wrapper.vm.formData = {
        Title: { value: "Original", action: "none" },
        Caption: { value: "Original", action: "none" },
      };
      expect(wrapper.vm.hasUnsavedChanges()).toBe(false);
    });

    it("should filter form data correctly", () => {
      const filtered = wrapper.vm.getFilteredFormData();

      expect(filtered).toEqual({
        Title: { action: "update", mixed: false, value: "Changed" },
      });
    });
  });

  describe("Location Functionality", () => {
    beforeEach(() => {
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Lat: { value: 37.7749, action: "none", mixed: false },
        Lng: { value: -122.4194, action: "none", mixed: false },
      };
      wrapper.vm.previousFormData = {
        Lat: { value: 40.7128 },
        Lng: { value: -74.006 },
      };
    });

    it("should handle location updates", () => {
      const newCoords = [40.7128, -74.006];
      wrapper.vm.updateLatLng(newCoords);

      expect(wrapper.vm.formData.Lat.value).toBe(40.7128);
      expect(wrapper.vm.formData.Lng.value).toBe(-74.006);
    });

    it("should handle location deletion", () => {
      wrapper.vm.onLocationDelete();

      expect(wrapper.vm.deletedFields.Lat).toBe(true);
      expect(wrapper.vm.deletedFields.Lng).toBe(true);
      expect(wrapper.vm.formData.Lat.value).toBe(0);
      expect(wrapper.vm.formData.Lng.value).toBe(0);
    });

    it("should handle location undo", () => {
      wrapper.vm.onLocationUndo();

      expect(wrapper.vm.deletedFields.Lat).toBe(false);
      expect(wrapper.vm.deletedFields.Lng).toBe(false);
      expect(wrapper.vm.formData.Lat.action).toBe("none");
      expect(wrapper.vm.formData.Lng.action).toBe("none");
    });

    it("should open location dialog", () => {
      wrapper.vm.adjustLocation();
      expect(wrapper.vm.locationDialog).toBe(true);
    });
  });

  describe("Save Functionality", () => {
    beforeEach(() => {
      wrapper.vm.model = mockBatch;
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Title: { value: "New Title", action: "update", mixed: false },
        Caption: { value: "New Caption", action: "update", mixed: false },
      };
    });

    it("should save changes successfully", async () => {
      await wrapper.vm.save(false);

      expect(mockBatch.save).toHaveBeenCalled();
      expect(global.$notify.success).toHaveBeenCalledWith("Changes successfully saved");
      expect(wrapper.vm.saving).toBe(false);
    });

    it("should handle save errors", async () => {
      mockBatch.save.mockRejectedValue(new Error("Save failed"));

      await wrapper.vm.save(false);

      expect(global.$notify.error).toHaveBeenCalledWith("Failed to save changes");
      expect(wrapper.vm.saving).toBe(false);
    });

  });

  describe("Form Field Updates", () => {
    beforeEach(() => {
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Title: { value: "Test", action: "none", mixed: false },
      };
      wrapper.vm.previousFormData = {
        Title: { value: "Original", action: "none" },
      };
    });

    it("should handle text field changes", () => {
      wrapper.vm.changeValue("New Title", "text-field", "Title");

      expect(wrapper.vm.formData.Title.value).toBe("New Title");
      expect(wrapper.vm.formData.Title.action).toBe("update");
    });

    it("should reset action when value returns to original", () => {
      wrapper.vm.changeValue("Original", "text-field", "Title");

      expect(wrapper.vm.formData.Title.value).toBe("Original");
      expect(wrapper.vm.formData.Title.action).toBe("none");
    });
  });

  describe("Selection Management", () => {
    beforeEach(() => {
      wrapper.vm.model = mockBatch;
    });

    it("should handle photo opening", () => {
      wrapper.vm.openPhoto(0);
      expect(global.$lightbox.openModels).toHaveBeenCalled();
    });
  });

  describe("Date Validation", () => {
    beforeEach(() => {
      wrapper.vm.formData = {
        ...wrapper.vm.formData,
        Year: { value: 2023, mixed: false },
        Month: { value: 2, mixed: false },
        Day: { value: 30, mixed: false, action: "update" },
      };
      wrapper.vm.actions = { update: "update", none: "none" };
    });

    it("should clamp day when date is resolvable", () => {
      wrapper.vm.clampBatchDayIfResolvable();

      // February 2023 has 28 days, so day should be clamped to 28
      expect(wrapper.vm.formData.Day.value).toBe(28);
      expect(wrapper.vm.formData.Day.action).toBe("update");
    });

    it("should not clamp when date is not resolvable", () => {
      wrapper.vm.formData.Year.mixed = true; // Make it non-resolvable

      wrapper.vm.clampBatchDayIfResolvable();

      // Should remain unchanged
      expect(wrapper.vm.formData.Day.value).toBe(30);
    });
  });

  describe("Component Lifecycle", () => {
    it("should initialize data when visible becomes true", async () => {
      await wrapper.setProps({ visible: true });
      await nextTick();
      await nextTick();
      expect(mockBatch.getData).toHaveBeenCalledWith(mockSelection);
    });

    it("should emit close event", () => {
      wrapper.vm.close();
      expect(wrapper.emitted("close")).toBeTruthy();
    });
  });
});
