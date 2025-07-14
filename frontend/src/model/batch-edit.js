import $api from "common/api";
import Model from "./model";
import { Photo } from "model/photo";

export class Batch extends Model {
  constructor(values) {
    super(values);
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
    return await $api.post("batch/photos/edit", { photos: selection }).then((response) => {
      const models = response.data.models;
      const modelsLength = response.data.models.length;

      if (modelsLength > 0) {
        for (let i = 0; i < modelsLength; i++) {
          const modelInstance = new Photo();
          this.models.push(modelInstance.setValues(models[i]));
        }
      }

      this.values = response.data.values;
      this.setSelections(selection);
    });
  }

  setSelections(selection) {
    this.selection = selection.map((id) => {
      return {
        id: id,
        selected: true,
      };
    });
  }

  isSelected(id) {
    let isSelected = null;

    this.selection.find((element) => {
      if (element.id === id) {
        isSelected = element.selected;
      }
    });

    return isSelected;
  }

  getLengthOfAllSelected() {
    return this.selection.filter((photo) => photo.selected).length;
  }

  toggle(id) {
    this.selection.find((element) => {
      if (element.id === id) {
        element.selected = !element.selected;
      }
    });
  }

  toggleAll(isToggledAll) {
    this.selection.find((element) => {
      element.selected = isToggledAll;
    });
  }
}
