<script>
import { h } from 'vue';

// Emojis (Icons as Vue components)

// Import static image
import sigenapls from '@/assets/images/gifts/sigenapls.png';
import isforu from '@/assets/images/gifts/isforu.png';
import isforyou from '@/assets/images/gifts/isforyou.png';
import justin from '@/assets/images/gifts/justin_timberlake_throwing_a_rose.jpg.png';
import meow from '@/assets/images/gifts/meow.png';
import me_literally from '@/assets/images/gifts/me_literally.png';
import rizzler from '@/assets/images/gifts/rizzler.png';
import sili from '@/assets/images/gifts/sili.png';
import batman from '@/assets/images/gifts/murag_batman.png';

const emojis = {
  sigenapls:  sigenapls,
  isforu: isforu,
  isforyou: isforyou,
  timberlake: justin,
  mukuha: meow,
  me_literally: me_literally,
  rizzler: rizzler,
  sili: sili,
  muragbatman : batman
};

export default {
  props: {
    uid: {
      type: String,
      required: true
    }
  },
  render() {
    const emoji = emojis[this.uid];

    // If it's a Vue component, render it
    if (typeof emoji === 'object') {
      return h(emoji, this.$props);
    }

    // If it's an image, render an <img> element
    return h('img', {
      src: emoji,
      alt: this.uid,
      ...this.$props
    });
  }
};
</script>
