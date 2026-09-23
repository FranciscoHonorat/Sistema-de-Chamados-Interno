import { onMounted, ref, type Ref } from 'vue'

export function useLoad<T>(
  load: () => Promise<T>,
  initial: T,
): { data: Ref<T>; failed: Ref<boolean>; loading: Ref<boolean>; reload: () => Promise<void> } {
  const data = ref(initial) as Ref<T>
  const failed = ref(false)
  const loading = ref(true)

  async function reload() {
    loading.value = true
    try {
      data.value = await load()
      failed.value = false
    } catch {
      failed.value = true
    } finally {
      loading.value = false
    }
  }

  onMounted(reload)

  return { data, failed, loading, reload }
}
