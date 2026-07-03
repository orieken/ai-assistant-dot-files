<!--
  Intentional test fixture for tests/platform-verification/ — NOT real framework code.
  Deliberately violates vue-frontend.mdc / vue-frontend.instructions.md rules:
    - Options API instead of Composition API with <script setup>
    - Custom CSS instead of Tailwind utility classes
    - No TypeScript
    - Business logic (the discount calculation) sitting directly in the component
    - Over 100 lines is not needed to demonstrate the violation, just the pattern itself
-->
<template>
  <div class="product-card">
    <h2>{{ name }}</h2>
    <p>{{ finalPrice }}</p>
    <button @click="addToCart">Add to cart</button>
  </div>
</template>

<script>
export default {
  name: 'ProductCard',
  props: ['name', 'price', 'isMember'],
  data() {
    return { inCart: false }
  },
  computed: {
    finalPrice() {
      // Business logic that belongs in a composable/service, not the component.
      if (this.isMember) {
        return this.price * 0.9 * 1.08
      }
      return this.price * 1.08
    }
  },
  methods: {
    addToCart() {
      this.inCart = true
    }
  }
}
</script>

<style>
.product-card {
  border: 1px solid #ccc;
  padding: 16px;
  border-radius: 8px;
}
</style>
