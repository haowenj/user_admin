<template>
  <div class="employee-container">
  <el-card class="employee-card" ref="cardRef">
      <template #header>
        <div class="header">
          <h3>员工信息管理</h3>
          <el-button type="primary" @click="handleAdd">新增员工</el-button>
        </div>
      </template>
      
      <el-table class="employee-table" :data="employees" border stripe :height="tableHeight" table-layout="auto">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="姓名" min-width="140" />
        <el-table-column prop="age" label="年龄" width="80" />
        <el-table-column prop="gender" label="性别" width="80" />
        <el-table-column prop="department" label="部门" min-width="160" />
        <el-table-column prop="position" label="职位" min-width="160" />
        <el-table-column prop="hire_date" label="入职日期" min-width="140" />
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑员工' : '新增员工'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="年龄" prop="age">
          <el-input-number v-model="form.age" :min="18" :max="65" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select v-model="form.gender">
            <el-option label="男" value="男" />
            <el-option label="女" value="女" />
          </el-select>
        </el-form-item>
        <el-form-item label="部门" prop="department">
          <el-input v-model="form.department" />
        </el-form-item>
        <el-form-item label="职位" prop="position">
          <el-input v-model="form.position" />
        </el-form-item>
        <el-form-item label="入职日期" prop="hire_date">
          <el-date-picker v-model="form.hire_date" type="date" value-format="YYYY-MM-DD" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { employeeAPI } from '../api'

const employees = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref(null)
const formRef = ref()
const cardRef = ref()
const tableHeight = ref(null)
let cardBodyEl = null
let resizeRaf = 0

const form = ref({
  name: '',
  age: 25,
  gender: '男',
  department: '',
  position: '',
  hire_date: ''
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  department: [{ required: true, message: '请输入部门', trigger: 'blur' }],
  position: [{ required: true, message: '请输入职位', trigger: 'blur' }]
}

const loadEmployees = async () => {
  try {
    const res = await employeeAPI.getList()
    employees.value = res.data
    await nextTick()
    updateTableHeight()
  } catch (err) {
    ElMessage.error('获取员工列表失败')
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.value = { name: '', age: 25, gender: '男', department: '', position: '', hire_date: '' }
  dialogVisible.value = true
}

const handleEdit = (row) => {
  isEdit.value = true
  currentId.value = row.id
  form.value = { ...row }
  dialogVisible.value = true
}

const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该员工吗？', '提示', { type: 'warning' })
    await employeeAPI.delete(row.id)
    ElMessage.success('删除成功')
    loadEmployees()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    if (isEdit.value) {
      await employeeAPI.update(currentId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await employeeAPI.create(form.value)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadEmployees()
  } catch (err) {
    ElMessage.error(err.response?.data?.error || '操作失败')
  }
}

const updateTableHeight = () => {
  if (!cardBodyEl) return
  const height = cardBodyEl.clientHeight
  if (height > 0) {
    tableHeight.value = height
  }
}

const handleResize = () => {
  if (resizeRaf) {
    cancelAnimationFrame(resizeRaf)
  }
  resizeRaf = requestAnimationFrame(() => {
    updateTableHeight()
  })
}

onMounted(async () => {
  await nextTick()
  const cardEl = cardRef.value?.$el
  cardBodyEl = cardEl ? cardEl.querySelector('.el-card__body') : null
  updateTableHeight()
  window.addEventListener('resize', handleResize)
  loadEmployees()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  if (resizeRaf) {
    cancelAnimationFrame(resizeRaf)
  }
})
</script>

<style scoped>
.employee-container {
  padding: 20px;
  height: 100%;
  box-sizing: border-box;
  min-height: 0;
}

.employee-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.employee-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.employee-table {
  flex: 1;
  width: 100%;
  height: 100%;
  min-height: 0;
}

h3 {
  margin: 0;
}
</style>
