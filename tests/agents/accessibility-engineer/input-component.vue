<!-- UserSearchBar.vue — deliberately contains accessibility issues for eval fixture -->
<template>
  <div class="search-bar">
    <input
      type="text"
      placeholder="Search users..."
      @input="onInput"
      @keydown.enter="onSubmit"
      style="color: #999; background: #eee;"
    />
    <div class="icon" @click="onSubmit">🔍</div>
    <div v-if="errorMsg" class="error" style="color: red;">
      {{ errorMsg }}
    </div>
    <ul>
      <li v-for="result in results" @click="selectUser(result)">
        {{ result.name }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ results: Array<{ id: string; name: string }> }>();
const emit = defineEmits(['select', 'search']);

const errorMsg = ref('');

function onInput(e: Event) {
  emit('search', (e.target as HTMLInputElement).value);
}

function onSubmit() {
  emit('search', '');
}

function selectUser(user: { id: string; name: string }) {
  emit('select', user);
}
</script>
