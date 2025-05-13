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
      // TODO: uncomment to test value
      // this.setFormValues();
    });
  }

  // TODO: uncomment to test value
  // setFormValues() {
  //   console.log('this.values', this.values);
  //
  //   Object.entries(this.values).forEach(([key, val]) => {
  //     if (val.mixed) {
  //       this.values[key] = "<mixed>";
  //     } else if (!val.mixed && !val.value) {
  //       this.values[key] = val.value;
  //     } else if (!val.mixed && val.value) {
  //       this.values[key] = 'Label';
  //     } else {
  //       this.values[key] = "";
  //     }
  //   });
  // }

  setSelections(selection) {
    this.selection = selection.map(id => {
      return {
        id: id,
        selected: true,
      };
    });
  }

  getPlaceholderForField(fieldType, fieldName) {
    const fieldData = this.values[fieldName];
    console.log('fieldName', fieldName);

    if(!fieldData) return;

    if(fieldType === 'text-field') {
      if (fieldData.mixed) {
        console.log("mixed", this.values[fieldName]);
        return "<mixed>";
      } else if (!fieldData.mixed && !fieldData.value) {
        console.log("!mixed && value", this.values[fieldName]);
        return fieldData.value;
      } else if (!fieldData.mixed && fieldData.value.isEmpty()) {
        console.log("!mixed && !value", this.values[fieldName]);
        return fieldName;
      } else {
        return "";
      }
    } else if(fieldType === 'input-field') {
      // TODO: change the logic
      return 'EMPTY';
    }
  }

  isSelected(id) {
    let isSelected = null;

    this.selection.find(element => {
      if (element.id === id) {
        isSelected = element.selected;
      }
    });

    return isSelected;
  }

  getLengthOfAllSelected() {
    return this.selection.filter(photo => photo.selected).length;
  }

  toggle(id) {
    this.selection.find(element => {
      if (element.id === id) {
        element.selected = !element.selected;
      }
    });
  }

  toggleAll(isToggledAll){
    this.selection.find(element => {
      element.selected = isToggledAll;
    });
  }
}
