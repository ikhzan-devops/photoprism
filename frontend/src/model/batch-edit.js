import $api from "common/api";
import Model from "./model";
import { Photo } from "model/photo";

export class Batch extends Model {
  constructor(values) {
    super(values);
    this.selectionById = new Map();
  }

  getDefaults() {
    return {
      models: [],
      values: {},
      selection: [],
    };
  }

  getDefaultFormData() {
    return {
      Title: {},
      DetailsSubject: {},
      Caption: {},
      Day: {},
      Month: {},
      Year: {},
      TimeZone: {},
      Country: {},
      Altitude: {},
      Lat: {},
      Lng: {},
      DetailsArtist: {},
      DetailsCopyright: {},
      DetailsLicense: {},
      DetailsKeywords: {},
      Type: {},
      Scan: {},
      Private: {},
      Favorite: {},
      Panorama: {},
      Iso: {},
      FocalLength: {},
      FNumber: {},
      Exposure: {},
      CameraID: {},
      LensID: {},
      Albums: {
        action: "none",
        mixed: false,
        items: [],
      },
      Labels: {
        action: "none",
        mixed: false,
        items: [],
      },
    };
  }

  save(selection, values) {
    return $api
      .post("batch/photos/edit", { photos: selection, values: values })
      .then((response) => {
        if (response.data.values) {
          this.values = response.data.values;
        }
        return Promise.resolve(this);
      })
      .catch((error) => {
        throw error;
      });
  }

  async getData(selection) {
    const response = await $api.post("batch/photos/edit", { photos: selection });
    const models = response.data.models || [];

    this.models = models.map((m) => {
      const modelInstance = new Photo();
      return modelInstance.setValues(m);
    });

    this.values = response.data.values;
    this.setSelections(selection);
  }

  async getValuesForSelection(selection) {
    const response = await $api.post("batch/photos/edit", { photos: selection });
    this.values = response.data.values;
    return this.values;
  }

  setSelections(selection) {
    this.selection = selection.map((id) => {
      return {
        id: id,
        selected: true,
      };
    });
    this.selectionById = new Map(this.selection.map((entry) => [entry.id, entry]));
  }

  isSelected(id) {
    const entry = this.selectionById && this.selectionById.get(id);
    return entry ? entry.selected : null;
  }

  getLengthOfAllSelected() {
    return this.selection.filter((photo) => photo.selected).length;
  }

  toggle(id) {
    const entry = this.selectionById && this.selectionById.get(id);
    if (entry) {
      entry.selected = !entry.selected;
    }
  }

  toggleAll(isToggledAll) {
    this.selection.forEach((element) => {
      element.selected = isToggledAll;
    });
  }

  wasChanged() {
    return super.wasChanged();
  }
}
