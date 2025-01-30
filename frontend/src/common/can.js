/*

Copyright (c) 2018 - 2025 PhotoPrism UG. All rights reserved.

    This program is free software: you can redistribute it and/or modify
    it under Version 3 of the GNU Affero General Public License (the "AGPL"):
    <https://docs.photoprism.app/license/agpl>

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    The AGPL is supplemented by our Trademark and Brand Guidelines,
    which describe how our Brand Assets may be used:
    <https://www.photoprism.app/trademark>

Feel free to send an email to hello@photoprism.app if you have questions,
want to support our work, or just want to say hello.

Additional information can be found in our Developer Guide:
<https://docs.photoprism.app/developer-guide/>

*/

import * as media from "common/media";

export const useVideo = !!document.createElement("video").canPlayType;
export const useAVC = useVideo // AVC
  ? !!document.createElement("video").canPlayType(media.ContentTypeAVC)
  : false;
export const useHEVC = useVideo // HEVC, Basic Support
  ? !!document.createElement("video").canPlayType(media.ContentTypeHEVC)
  : false;
export const useHEV1 = useVideo // HEV1, Basic Support
  ? !!document.createElement("video").canPlayType(media.ContentTypeHEV1)
  : false;
export const useVVC = useVideo // VVC, Basic Support
  ? !!document.createElement("video").canPlayType(media.ContentTypeVVC)
  : false;
export const useOGV = useVideo // Ogg Theora
  ? !!document.createElement("video").canPlayType(media.ContentTypeOGV)
  : false;
export const useWebM = useVideo // Google WebM
  ? !!document.createElement("video").canPlayType(media.ContentTypeWebM)
  : false;
export const useVP8 = useVideo // Google WebM, VP8
  ? !!document.createElement("video").canPlayType(media.ContentTypeVP8)
  : false;
export const useVP9 = useVideo // Google WebM, VP9
  ? !!document.createElement("video").canPlayType(media.ContentTypeVP9)
  : false;
export const useAV1 = useVideo // AV1, Main Profile
  ? !!document.createElement("video").canPlayType(media.ContentTypeAV1)
  : false;
