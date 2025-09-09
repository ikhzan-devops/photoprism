import { Selector, t } from "testcafe";

export default class Page {
  constructor() {}

  async openContextMenu() {
    if (!(await Selector(".action-clear").visible)) {
      await t.click(Selector("button.action-menu"));
    }
  }

  async checkContextMenuCount(count) {
    const Count = await Selector("span.count-clipboard", { timeout: 5000 });
    await t.expect(Count.textContent).eql(count);
  }

  async checkContextMenuActionAvailability(action, visible) {
    await this.openContextMenu();
    if (visible) {
      await t
        .expect(Selector("#t-clipboard button.action-" + action).visible)
        .ok()
        .expect(Selector("#t-clipboard button.action-" + action).hasAttribute("disabled"))
        .notOk();
    } else {
      if (await Selector("#t-clipboard button.action-" + action).visible) {
        await t.expect(Selector("#t-clipboard button.action-" + action).hasAttribute("disabled")).ok();
      } else {
        await t.expect(Selector("#t-clipboard button.action-" + action).visible).notOk();
      }
    }
  }
  async triggerContextMenuAction(action, albumName) {
    await this.openContextMenu();
    if (t.browser.platform === "mobile") {
      await t.wait(5000);
    }
    await t.click(Selector("#t-clipboard button.action-" + action));
    if (action === "delete") {
      await t.click(Selector("button.action-confirm"));
    }
    if ((action === "album") | (action === "clone")) {
      await t.click(Selector(".input-albums"));

      // Handle single album name or array of album names
      const albumNames = Array.isArray(albumName) ? albumName : [albumName];

      for (const name of albumNames) {
        if (await Selector("div").withText(name).parent('div[role="option"]').visible) {
          // Click on the album option to select it
          await t.click(Selector("div").withText(name).parent('div[role="option"]'));
        } else {
          // Type the new album name and press enter to create it
          await t.typeText(Selector(".input-albums input"), name).pressKey("enter");
        }

        // Wait a bit for the UI to update after selection
        await t.wait(500);
      }

      await t.click(Selector("button.action-confirm"));
    }
  }

  async clearSelection() {
    await this.openContextMenu();
    await t.click(Selector(".action-clear"));
  }

  async triggerAlbumDialogAndType(albumName) {
    await this.openContextMenu();
    if (t.browser.platform === "mobile") {
      await t.wait(5000);
    }
    await t.click(Selector("#t-clipboard button.action-album"));

    // Wait for dialog to open
    await t.wait(1000);

    // Type the album name and press enter (simulating the bug scenario)
    await t
      .click(Selector(".input-albums input"))
      .typeText(Selector(".input-albums input"), albumName, { replace: true })
      .pressKey("enter");

    // Return the dialog selectors for further testing
    return {
      confirmButton: Selector("button.action-confirm"),
      cancelButton: Selector("button").withText("Cancel"),
      inputField: Selector(".input-albums input"),
      chips: Selector("v-chip"),
    };
  }
}
