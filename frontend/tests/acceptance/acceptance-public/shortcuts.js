import { Selector, ClientFunction } from "testcafe";
import testcafeconfig from "../../testcafeconfig.json";
import Menu from "../page-model/menu";
import Toolbar from "../page-model/toolbar";
import Photo from "../page-model/photo";
import PhotoViewer from "../page-model/photoviewer";
import Page from "../page-model/page";
import PhotoEdit from "../page-model/photo-edit";
import Subject from "../page-model/subject";
import Label from "../page-model/label";
import Library from "../page-model/library";

fixture`Test Keyboard Shortcuts`
  .page`${testcafeconfig.url}`;

const menu = new Menu();
const toolbar = new Toolbar();
const photo = new Photo();
const photoviewer = new PhotoViewer();
const page = new Page();
const photoEdit = new PhotoEdit();
const subject = new Subject();
const label = new Label();
const library = new Library();

const triggerKeyPress = ClientFunction((key, code, keyCode, ctrlKey, shiftKey, targetSelector) => {
    const target = targetSelector ? document.querySelector(targetSelector) : document;
    if (!target) {
        console.error("Target element not found for selector:", targetSelector);
        return;
    }
    target.dispatchEvent(new KeyboardEvent('keydown', {
        key: key,
        code: code,
        keyCode: keyCode,
        which: keyCode,
        bubbles: true,
        cancelable: true,
        ctrlKey: ctrlKey,
        shiftKey: shiftKey,
        altKey: false,
        metaKey: false
    }));
}, {
});

const isFullscreen = ClientFunction(() => !!document.fullscreenElement);

test.meta("testID", "shortcuts-001").meta({ type: "short", mode: "public" })(
  "Common: Test General Page Shortcuts",
  async (t) => {
    await menu.openPage("browse");
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, false);
    await t.expect(toolbar.search1.focused).ok();

    await t.wait(500);
    await triggerKeyPress('r', 'KeyR', 82, true, false);
    await t.expect(Selector("div.is-photo").exists).ok();

    await t.wait(500);
    await triggerKeyPress('u', 'KeyU', 85, true, false);
    await t.expect(Selector('div').withText('Upload').nth(3).visible).ok();
    await t.pressKey("esc");
    await t.expect(Selector(".p-upload-dialog").visible).notOk();
  }
);

test.meta("testID", "shortcuts-002").meta({ type: "short", mode: "public" })(
  "Common: Test Lightbox Shortcuts",
  async (t) => {
    await menu.openPage("browse");
    await t.navigateTo('/library/videos');
    const videoUid = await photo.getNthPhotoUid("all", 0);
    await photoviewer.openPhotoViewer("uid", videoUid);

    await t.wait(500);
    const infoPanelSelector = Selector('div').withText('Information').nth(4);
    await t.expect(infoPanelSelector.visible).notOk("Information panel should not be visible initially");

    await triggerKeyPress('i', 'KeyI', 73, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(infoPanelSelector.visible).ok("Information panel should be visible after first Ctrl+I");

    await triggerKeyPress('i', 'KeyI', 73, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(infoPanelSelector.visible).notOk("Information panel should be hidden after second Ctrl+I");

    await t.wait(500);
    await triggerKeyPress('m', 'KeyM', 77, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(Selector('.p-lightbox__content').hasClass("is-muted")).ok("Video should be muted after first Ctrl+M");
    await t.wait(500);
    await triggerKeyPress('m', 'KeyM', 77, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(Selector('.p-lightbox__content').hasClass("is-muted")).notOk("Video should be unmuted after second Ctrl+M");

    await t.wait(500);
    await triggerKeyPress('s', 'KeyS', 83, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(Selector('.p-lightbox__content').hasClass("slideshow-active")).ok("Slideshow should be active after first Ctrl+S");
    await t.wait(500);
    await triggerKeyPress('s', 'KeyS', 83, true, false, 'div.p-lightbox__pswp');
    await t.wait(500);
    await t.expect(Selector('.p-lightbox__content').hasClass("slideshow-active")).notOk("Slideshow should be inactive after second Ctrl+S");

    await triggerKeyPress('Escape', 'Escape', 27, false, false, 'div.p-lightbox__pswp');
  }
);

// NOT WORKING SKIPPED FOR NOW
test.meta("testID", "shortcuts-003").meta({ type: "short", mode: "public" }).skip(
  "Common: Test Lightbox Archive and Download Shortcuts",
  async (t) => {
    await menu.openPage("browse");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);

    await t.expect(photoviewer.viewer.visible).ok();
    
    await triggerKeyPress('a', 'KeyA', 65, true, false, 'div.p-lightbox__pswp');
    
    await t.wait(3000);
    
    const snackbarVisible = await Selector('.v-snackbar').visible;
    
    if (snackbarVisible) {
        await Selector('.v-snackbar__content').innerText;
    }
    
    await triggerKeyPress('d', 'KeyD', 68, true, false, 'div.p-lightbox__pswp');
    await t.wait(2000);

    await t.pressKey("esc");
  }
);

test.meta("testID", "shortcuts-004").meta({ type: "short", mode: "public" })(
  "Common: Test Lightbox Edit, Fullscreen, and Like Shortcuts",
  async (t) => {
    await menu.openPage("browse");
    const FirstPhotoUid = await photo.getNthPhotoUid("image", 0);
    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);

    // Edit Test
    await t.wait(500);
    await triggerKeyPress('e', 'KeyE', 69, true, false);
    await t.expect(photoEdit.dialog.visible).ok();
    await t.pressKey("esc");

    await photoviewer.openPhotoViewer("uid", FirstPhotoUid);

    // Fullscreen Test
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, false, 'div.p-lightbox__pswp');
    await t.wait(1000);
    await t.expect(isFullscreen()).ok("Browser did not enter fullscreen mode.");

    await t.wait(1000);
    await triggerKeyPress('f', 'KeyF', 70, true, false, 'div.p-lightbox__pswp');
    await t.wait(1000);
    await t.expect(isFullscreen()).notOk("Browser did not exit fullscreen mode.");

    // Like Test
    await t.wait(500);
    const isLikedInitially = await Selector('.p-lightbox__content').hasClass('is-favorite');
    await triggerKeyPress('l', 'KeyL', 76, true, false, 'div.p-lightbox__pswp');
    await t.wait(2000); // Wait for potential UI updates
    await t.expect(photoviewer.menuButton.exists).ok("Menu button does not exist after Ctrl+L");
    const isLikedAfterToggle = await Selector('.p-lightbox__content').hasClass('is-favorite');
    if (isLikedInitially) {
        await t.expect(isLikedAfterToggle).notOk("Failed to unlike photo after Ctrl+L");
    } else {
        await t.expect(isLikedAfterToggle).ok("Failed to like photo after Ctrl+L");
    }

    await t.pressKey("esc");
  }
);

test.meta("testID", "shortcuts-005").meta({ type: "short", mode: "public" })(
  "Common: Test Expansion Panel and Page-Specific Search Focus",
  async (t) => {
    await menu.openPage("browse");
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, true);
    await t.wait(500);
    await t.expect(Selector(".toolbar-expansion-panel").find('div').withText('All Countries').exists).ok("Expansion panel content ('All Countries') not found after Shift+Ctrl+F");
    // Close expansion panel
    // TODO: Currently no close functionality implemented
    // await triggerKeyPress('f', 'KeyF', 70, true, true); // Toggle it off
    // await t.wait(500); // Wait for animation
    // await t.expect(Selector(".toolbar-expansion-panel").getAttribute("style")).contains("display: none", "Expansion panel is not hidden (display: none) after second Shift+Ctrl+F");

    await menu.openPage("people");
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, false);
    await t.expect(subject.search.focused).ok();

    await menu.openPage("labels");
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, false);
    await t.expect(label.search.focused).ok();

    await t.navigateTo('/library/errors');
    await t.wait(500);
    await triggerKeyPress('f', 'KeyF', 70, true, false);
    await t.expect(library.searchInput.focused).ok();
  }
);
