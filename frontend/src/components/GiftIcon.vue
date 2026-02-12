<script>
import { h } from 'vue';

// Emoji mappings for gift UIDs
const emojiMap = {
  // Standard Valentine gifts
  rose: '🌹',
  chocolate: '🍫',
  teddy: '🧸',
  flowers: '💐',
  candy: '🍬',
  card: '💌',
  balloon: '🎈',
  money_100: '💵',
  money_500: '💰',
  
  // Legacy/custom gifts (keep existing)
  sigenapls: '🌹',
  isforu: '💝',
  isforyou: '💝',
  timberlake: '🌹',
  mukuha: '😺',
  me_literally: '🤷',
  rizzler: '😎',
  sili: '🌶️',
  muragbatman: '🦇'
};

export default {
  props: {
    uid: {
      type: String,
      required: true
    }
  },
  render() {
    // Clean the UID (remove whitespace/tabs that might exist in database)
    const cleanUid = this.uid?.trim() || '';
    const emoji = emojiMap[cleanUid] || '🎁'; // Default gift emoji

    // Render emoji as text
    return h('span', {
      ...this.$props,
      style: 'font-size: inherit; line-height: 1; display: inline-block;'
    }, emoji);
  }
};
</script>
