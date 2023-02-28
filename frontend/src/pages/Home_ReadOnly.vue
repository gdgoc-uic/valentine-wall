<template>
  <main class="max-w-2xl mx-auto w-full px-2 pt-8 pb-2">
    <img src="../assets/images/logo.png" class="w-4/5 md:w-2/3 lg:w-2/3 pb-8 mx-auto" alt="Valentine Wall">
    <div class="bg-white rounded-lg shadow-lg p-8 flex flex-col items-center">
      <div class="prose text-center">
        <p>Since the website's launch, we have received <b>{{ messagesCount }}</b> messages that ranges from simple appreciations to Valorant invites to words of wisdom all of us can relate to!</p>
        <p>We closed the website last {{ closingDate }} and many of you weren't able to get a chance to read or view them all. Thus, we are giving everyone the chance to archive their messages until <b>{{ shuttingDownDate }}</b>.</p>
        <p v-if="!authState.isLoggedIn">To start archiving, please click the sign-in button below.</p>
        <p v-else>To start archiving, please click the "Archive" below.</p>
      </div>

      <div class="flex flex-col space-y-2 my-8 ">
        <login-button v-if="!authState.isLoggedIn" class="btn-lg lg:px-24" />
        <template v-else>
          <archive-dialog @error="catchAndNotifyError" v-slot="{ openDialog }">
            <button @click="openDialog" class="btn btn-primary btn-lg lg:px-24">Archive</button>
          </archive-dialog>
          <button @click="logout()" class="btn btn-ghost">Log out</button>
        </template>
      </div>

      <div class="flex flex-col space-y-4 pt-4">
        <div :key="'faq_' + fi" v-for="(faq, fi) in faqs" class="collapse border border-base-300 bg-base-100 rounded-box collapse-plus">
          <input type="checkbox">
          <div class="collapse-title text-xl font-medium">{{ faq.title }}</div>
          <div v-html="faq.content" class="collapse-content prose-lg prose-a:underline prose-a:text-rose-500"></div>
        </div>
      </div>
    </div>
  </main>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { pb } from '../client';
import ArchiveDialog from '../components/ArchiveDialog.vue';
import LoginButton from '../components/LoginButton.vue';
import { catchAndNotifyError } from '../notify';
import { useAuth } from '../store_new';

const { state: authState, methods: { logout } } = useAuth();
const messagesCount = ref(0);
const closingDate = "March 1";
const shuttingDownDate = "March 14, 2023";

async function loadStats(): Promise<void> {
  try {
    const giftsResp = await pb.collection('messages').getList(1, 1);
    messagesCount.value = giftsResp.totalItems;
  } catch(e) {
    catchAndNotifyError(e);
  }
}

const faqs = [
  {
    title: 'How to get my messages?',
    content: `
      <p>Simply sign-in to your UIC Google account by clicking the big button you're seeing above the screen. Once you have logged in, an archive button will appear.</p>
    `
  },
  {
    title: 'What\'s included in the archive?',
    content: `
      <p>The downloaded ZIP archive contains the image versions of all the messages your student ID has received as well as a summary.</p>
    `
  },
  {
    title: 'Does it include the list of people who sent those messages?',
    content: `
      <p>Unfortunately, no.</p>
    `
  },
  {
    title: 'Can unregistered accounts still download messages?',
    content: `
      <p>Yes! As long as there are messages received to that email's student ID, you can download them.</p>
    `
  },
  {
    title: 'Can I download other people\'s messages?',
    content: `
      <p>No you cannot. That's why we required everyone to use the sign-in button to avoid any accidental leaks.</p>
    `
  },
  {
    title: 'Do you have any plans to post this to your Facebook page?',
    content: `
      <p>Yes! We will be posting all the public messages to our Facebook page. If you have a public message that you have sent in which you would like to keep it as a secret, please message us at our <a href="https://facebook.com/dscuic">Facebook page</a>.</p>
    `
  },
  {
    title: `What happens after ${shuttingDownDate}?`,
    content: `
      <p>We will be shutting off this website for good and delete all the data stored so grab this opportunity to get your messages!</p>
    `
  }
];

onMounted(() => {
  loadStats();
});
</script>