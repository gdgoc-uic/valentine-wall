<html>
  <head>
    <title>image</title>
    <style>
      /* lato-regular - latin */
      @font-face {
        font-family: 'Lato';
        font-style: normal;
        font-weight: 400;
        src: url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.eot'); /* IE9 Compat Modes */
        src: local(''),
            url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.eot?#iefix') format('embedded-opentype'), /* IE6-IE8 */
            url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.woff2') format('woff2'), /* Super Modern Browsers */
            url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.woff') format('woff'), /* Modern Browsers */
            url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.ttf') format('truetype'), /* Safari, Android, iOS */
            url('{{ .BackendURL }}/renderer_assets/fonts/lato-v22-latin-regular.svg#Lato') format('svg'); /* Legacy iOS */
      }

      /* nanum-pen-script-regular - latin */
      @font-face {
        font-family: 'Nanum Pen Script';
        font-style: normal;
        font-weight: 400;
        src: url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.eot'); /* IE9 Compat Modes */
        src: local(''),
            url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.eot?#iefix') format('embedded-opentype'), /* IE6-IE8 */
            url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.woff2') format('woff2'), /* Super Modern Browsers */
            url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.woff') format('woff'), /* Modern Browsers */
            url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.ttf') format('truetype'), /* Safari, Android, iOS */
            url('{{ .BackendURL }}/renderer_assets/fonts/nanum-pen-script-v15-latin-regular.svg#NanumPenScript') format('svg'); /* Legacy iOS */
      }

      * {
        box-sizing: border-box;
      }
      body {
        margin: 0;
        font-size: 16px;
      }
      p {
        margin: 0;
      }
      .image-wrapper {
        width: 1200px;
        height: 675px;
        background-image: url({{ .BackendURL }}/renderer_assets/images/background.png);
        background-size: cover;
        padding: 3rem 6rem;
        position: relative;
        overflow: hidden;
      }
      .image-wrapper .content-wrapper {
        background: linear-gradient(to bottom,rgb(254, 243, 199) 3.35rem,#00b0d7 1px); background-size: 100% 3.5rem;
        height: 100%;
        border-radius: 2rem;
        padding-left:3rem;
        padding-right: 3rem;
        padding-top: 4rem;
        padding-bottom: 1.5rem;
        display: flex;
        flex-direction: column;
        text-overflow: ellipsis;
        box-shadow: 1px 4px 8px 0px rgba(244, 63, 94,0.43);
      }

      .image-wrapper .content-wrapper .content {
        font-family: 'Nanum Pen Script', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        line-height: 0.9;
        text-align: left;
        text-overflow: ellipsis;
        overflow: hidden;
        width: 100%;
        font-size: 3.5rem;
        line-height: 1;
        height: 100%;
      }
      .image-wrapper .content-wrapper .timestamp {
        font-family: 'Lato', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: rgb(10, 10, 10);
        position: relative;
        top: 1rem;
      }

      .logo {
        position: absolute;
        bottom: 1%;
        left: 50%;
        transform: translateX(-50%);
        width: 30%;
      }

      .gifts {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
      }

      .gift-0, .gift-2, .gift-1 {
        width: 13%;
        position: absolute;
        transform-origin: bottom center;
        filter: drop-shadow(1px 4px 8px rgba(128, 6, 42, 0.678));
      }

      .gift-0 {
        bottom: 5%;
        right: 5%;
        transform: rotate(30deg);
      }

      .gift-2 {
        top: 1%;
        left: 0;
        transform: rotate(20deg);
      }

      .gift-1 {
        top: 1%;
        right: 3%;
        transform: rotate(-20deg);
      }
    </style>
  </head>
  <body>
    <div id="image-preview" class="image-wrapper">
      <div class="content-wrapper">
        <p class="content">{{ .Content }}</p>
        <p class="timestamp">Posted on {{ .CreatedAt }}</p>
      </div>
      <div class="gifts">
        {{ range $i, $giftId := .GiftIDs }}
          <img src="./emojis/{{ $giftId }}.svg" class="gift-{{ $i }}" />
        {{ end }}
      </div>
      <img class="logo" src="{{ .BackendURL }}/renderer_assets/images/logo.png" />
    </div>
  </body>
</html>