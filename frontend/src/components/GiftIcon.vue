<script>
import { h } from 'vue';

// Emojis
import IconSunflower from '~icons/twemoji/sunflower';
import IconRose from '~icons/twemoji/rose';
import IconBalloon from '~icons/twemoji/balloon';
import IconTeddyBear from '~icons/twemoji/teddy-bear';
import IconRing from '~icons/twemoji/ring';
import IconMoney from '~icons/twemoji/money-with-wings';
import IconHeart from '~icons/twemoji/heart-with-ribbon';
import IconChocolate from '~icons/twemoji/chocolate-bar';
import IconPizza from '~icons/twemoji/pizza';
import IconBubbleTea from '~icons/twemoji/bubble-tea';

const emojis = {
  'sunflower': IconSunflower,
  'rose': IconRose,
  'balloons': IconBalloon,
  'teddy-bear': IconTeddyBear,
  'ring': IconRing,
  'money': IconMoney,
  'heart': IconHeart,
  'chocolate': IconChocolate,
  'pizza': IconPizza,
  'milk-tea': IconBubbleTea
};

export default {
  render() {
    return h(emojis[this.uid], this.$props);
  },
  props: {
    uid: {
      type: String,
      required: true
    }
  }
}
</script>