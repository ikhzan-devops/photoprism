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
      Albums: [],
      Labels: [],
      Scan: {},
      Private: {},
      Favorite: {},
      Panorama: {},
    };
  }

  save() {
    return $api
      .post("batch/photos/edit", this.getValues(true))
      .then((response) => Promise.resolve(this.setValues(response.data)));
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
