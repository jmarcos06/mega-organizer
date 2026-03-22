<template>
  <div>
    <!-- Navbar -->
    <nav class="bg-slate-800 border-b border-slate-700 p-4 sticky top-0 z-50">
      <div class="container mx-auto flex justify-between items-center">
        <div class="flex items-center gap-2">
          <i class="fas fa-clover text-green-400 text-2xl"></i>
          <h1 class="text-xl font-bold tracking-wider">MEGA HUB</h1>
        </div>
        <div class="text-sm text-slate-400">
          <span class="inline-block w-2 h-2 rounded-full bg-green-500 mr-1 animate-pulse"></span>
          Ao Vivo
        </div>
      </div>
    </nav>

    <div class="container mx-auto p-4 max-w-6xl grid grid-cols-1 lg:grid-cols-3 gap-6 mt-6">
      
      <!-- Left Column -->
      <div class="lg:col-span-1 space-y-6">
        <!-- Generator Form -->
        <div class="bg-slate-800 rounded-xl p-6 shadow-lg border border-slate-700">
          <h2 class="text-lg font-semibold mb-4 text-green-400"><i class="fas fa-plus-circle"></i> Novo Jogo</h2>
          <form @submit.prevent="submitBet" class="space-y-4">
            <div>
              <label class="block text-xs uppercase text-slate-400 mb-1">Temporada</label>
              <select v-model="form.season" class="w-full bg-slate-900 border border-slate-700 rounded p-2 focus:border-green-500 focus:outline-none text-white">
                <option v-for="s in availableSeasons" :key="s" :value="s">{{ s }}</option>
              </select>
            </div>
            <div>
              <label class="block text-xs uppercase text-slate-400 mb-1">Seu Nickname</label>
              <input v-model="form.nickname" type="text" required placeholder="Ex: DanielBoludo" class="w-full bg-slate-900 border border-slate-700 rounded p-2 focus:border-green-500 focus:outline-none text-white">
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs uppercase text-slate-400 mb-1">Dezenas (6-20)</label>
                <input v-model.number="form.qtd" @input="updateCost" type="number" min="6" max="20" class="w-full bg-slate-900 border border-slate-700 rounded p-2 focus:border-green-500 focus:outline-none">
              </div>
              <div>
                <label class="block text-xs uppercase text-slate-400 mb-1">Custo Est.</label>
                <div class="py-2 text-green-400 font-bold">R$ {{ formatMoney(estimatedCost) }}</div>
              </div>
            </div>
            <div>
               <label class="block text-xs uppercase text-slate-400 mb-1">Números Fixos (Opcional)</label>
               <input v-model="form.fixos" type="text" placeholder="Ex: 7, 13, 50" class="w-full bg-slate-900 border border-slate-700 rounded p-2 focus:border-green-500 focus:outline-none text-sm">
               <p class="text-xs text-slate-500 mt-1">Separe por vírgula ou espaço.</p>
            </div>
            <button :disabled="isSubmitting" type="submit" class="w-full bg-green-600 hover:bg-green-500 text-white font-bold py-3 rounded transition-all transform hover:scale-105 disabled:opacity-50">
              <i class="fas fa-dice" v-if="!isSubmitting"></i> 
              <i class="fas fa-spinner fa-spin" v-else></i> 
              {{ isSubmitting ? 'Processando...' : 'Gerar e Apostar' }}
            </button>
          </form>
        </div>

        <!-- Dashboard Stats -->
        <div class="bg-slate-800 rounded-xl p-6 shadow-lg border border-slate-700">
          <h3 class="text-sm font-semibold text-slate-400 mb-3 uppercase">Estatísticas do Grupo</h3>
          <div class="flex justify-between items-center mb-2">
            <span>Total Apostado:</span>
            <span class="text-xl font-bold text-green-400">R$ {{ formatMoney(stats.total_gasto) }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span>Jogos Feitos:</span>
            <span class="text-xl font-bold text-blue-400">{{ stats.total_jogos }}</span>
          </div>
          <div class="mt-4 pt-4 border-t border-slate-700">
            <p class="text-xs text-slate-500 mb-2">NÚMEROS MAIS JOGADOS:</p>
            <div class="flex flex-wrap gap-2">
              <span v-for="n in stats.numeros_quentes" :key="n.numero" class="bg-yellow-600 text-yellow-100 text-xs px-2 py-1 rounded-full font-bold">
                {{ n.numero }} ({{ n.qtd }}x)
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Feed -->
      <div class="lg:col-span-2">
        <div class="bg-slate-800 rounded-xl shadow-lg border border-slate-700 flex flex-col h-[600px]">
          <div class="p-4 border-b border-slate-700 flex justify-between items-center">
            <h2 class="font-semibold text-slate-200"><i class="fas fa-stream"></i> Apostas Recentes</h2>
            <div class="flex gap-2 items-center">
              <select v-model="currentSeason" @change="fetchData" class="bg-slate-900 border border-slate-700 rounded px-2 py-1 text-xs text-white">
                <option value="">Todas as temporadas</option>
                <option v-for="s in allSeasons" :key="s" :value="s">{{ s }}</option>
              </select>
              <button @click="fetchData" class="text-slate-400 hover:text-white"><i class="fas fa-sync-alt"></i></button>
            </div>
          </div>
          
          <div class="flex-1 overflow-y-auto p-4 space-y-3 custom-scrollbar">
            <div v-if="!stats.ultimas_apostas || stats.ultimas_apostas.length === 0" class="text-center text-slate-500 mt-10">Nenhuma aposta ainda. Seja o primeiro!</div>
            
            <div v-for="aposta in stats.ultimas_apostas" :key="aposta.id" class="bg-slate-750 hover:bg-slate-700 transition p-4 rounded-lg border border-slate-700/50">
              <div class="flex justify-between items-start mb-2">
                <div>
                  <span @click="showUserHistory(aposta.nickname)" class="font-bold text-white text-lg mr-2 cursor-pointer hover:text-green-400">{{ aposta.nickname }}</span>
                  <span :class="['text-xs px-2 py-0.5 rounded border', aposta.tipo === 'Simples' ? 'bg-green-900 text-green-300 border-green-700' : 'bg-purple-900 text-purple-300 border-purple-700']">{{ aposta.tipo }}</span>
                  <span class="text-xs bg-blue-900 text-blue-300 px-1 py-0.5 rounded border border-blue-700 ml-1">{{ aposta.season }}</span>
                </div>
                <div class="text-right text-xs text-slate-400 flex items-center">
                  <div>
                    <div>{{ new Date(aposta.data).toLocaleTimeString() }}</div>
                    <div class="text-slate-500 text-[10px]">{{ aposta.id }}</div>
                  </div>
                  <button @click="deleteAposta(aposta.id)" class="text-red-400 hover:text-red-300 text-xs ml-2" title="Deletar aposta"><i class="fas fa-trash"></i></button>
                </div>
              </div>
              <div class="flex flex-wrap gap-2 mb-2">
                <span v-for="n in aposta.numeros" :key="n" class="inline-block w-8 h-8 leading-8 text-center rounded-full bg-slate-700 text-slate-200 font-bold text-sm border border-slate-600 shadow-sm">{{ n }}</span>
              </div>
              <div class="text-xs text-slate-500 flex justify-between">
                <span>Custo: R$ {{ formatMoney(aposta.custo) }}</span>
                <button @click="reuseNumbers(aposta.numeros)" class="text-blue-400 hover:text-blue-300"><i class="fas fa-copy"></i> Reusar números</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- User Modal -->
    <div v-if="isModalOpen" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4" @click.self="isModalOpen = false">
      <div class="bg-slate-800 rounded-xl max-w-4xl w-full max-h-[90vh] flex flex-col">
        <div class="p-4 border-b border-slate-700 flex justify-between items-center">
          <h2 class="text-lg font-bold text-white">Histórico de {{ selectedNickname }}</h2>
          <button @click="isModalOpen = false" class="text-slate-400 hover:text-white"><i class="fas fa-times text-xl"></i></button>
        </div>
        <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
          <div v-if="userHistory.length === 0" class="text-center text-slate-500">Este usuário ainda não fez apostas.</div>
          <div v-else>
            <div v-for="(apostas, season) in getHistoryBySeason()" :key="season" class="mb-6">
              <h3 class="text-lg font-bold text-green-400 mb-3 border-b border-slate-600 pb-2">{{ season }}</h3>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div v-for="aposta in apostas" :key="aposta.id" class="bg-slate-700 p-3 rounded border border-slate-600">
                  <div class="flex justify-between items-center mb-2">
                    <span class="text-xs text-slate-400">{{ new Date(aposta.data).toLocaleString() }}</span>
                    <span :class="['text-xs', aposta.tipo === 'Simples' ? 'text-green-400' : 'text-purple-400']">{{ aposta.tipo }}</span>
                  </div>
                  <div class="mb-2">
                    <span v-for="n in aposta.numeros" :key="n" class="inline-block w-6 h-6 leading-6 text-center rounded-full bg-slate-600 text-slate-200 font-bold text-xs border border-slate-500 mr-1 mb-1">{{ n }}</span>
                  </div>
                  <div class="flex justify-between items-center text-xs">
                    <span class="text-slate-500">R$ {{ formatMoney(aposta.custo) }}</span>
                    <button @click="reuseNumbers(aposta.numeros); isModalOpen = false;" class="text-blue-400 hover:text-blue-300"><i class="fas fa-copy"></i> Reusar</button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const form = ref({
  season: 'Mega da Virada 2025',
  nickname: '',
  qtd: 6,
  fixos: ''
});
const estimatedCost = ref(6.00);
const isSubmitting = ref(false);

const stats = ref({
  total_gasto: 0,
  total_jogos: 0,
  ultimas_apostas: [],
  numeros_quentes: []
});
const allSeasons = ref([]);
const currentSeason = ref('');

const availableSeasons = ref(['Mega da Virada 2025', 'Mega Sena Regular', 'Verão 2025']);

const isModalOpen = ref(false);
const selectedNickname = ref('');
const userHistory = ref([]);
let intervalId = null;

const formatMoney = (val) => Number(val || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 });

const updateCost = async () => {
  if (form.value.qtd < 6 || form.value.qtd > 20) return;
  try {
    const res = await fetch(`/api/custo?qtd=${form.value.qtd}`);
    const data = await res.json();
    if(data.custo !== undefined) estimatedCost.value = data.custo;
  } catch(e) { console.error(e); }
};

const fetchData = async () => {
  try {
    const p = currentSeason.value ? `?season=${encodeURIComponent(currentSeason.value)}` : '';
    const res = await fetch(`/api/dados${p}`);
    const data = await res.json();
    stats.value = data;
    if (data.seasons && JSON.stringify(data.seasons) !== JSON.stringify(allSeasons.value)) {
      allSeasons.value = data.seasons;
      data.seasons.forEach(s => {
        if (!availableSeasons.value.includes(s)) availableSeasons.value.push(s);
      });
    }
  } catch(e) { console.error(e); }
};

const submitBet = async () => {
  isSubmitting.value = true;
  try {
    const res = await fetch('/api/apostar', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || 'Erro ao realizar aposta');

    if (data.colisao) {
      alert(`⚠️ ATENÇÃO: O usuário ${data.colisao_com} já fez esse jogo exato!`);
    } else {
      if (window.confetti) window.confetti({ particleCount: 100, spread: 70, origin: { y: 0.6 } });
    }

    form.value.fixos = '';
    form.value.qtd = 6;
    form.value.season = 'Mega da Virada 2025';
    updateCost();
    fetchData();
  } catch(e) {
    alert('Erro: ' + e.message);
  } finally {
    isSubmitting.value = false;
  }
};

const deleteAposta = async (id) => {
  if (!confirm('Tem certeza que deseja deletar esta aposta?')) return;
  const token = prompt('Digite o token de acesso para deletar:');
  if (!token) return;
  
  try {
    const res = await fetch(`/api/aposta/deletar?id=${id}&token=${token}`, { method: 'DELETE' });
    if (!res.ok) {
      const e = await res.json();
      throw new Error(e.error || 'Erro ao deletar');
    }
    fetchData();
    alert('Aposta deletada com sucesso!');
  } catch(e) {
    alert('Erro ao deletar: ' + e.message);
  }
};

const showUserHistory = async (nickname) => {
  selectedNickname.value = nickname;
  try {
    const res = await fetch(`/api/usuario/historico?nickname=${encodeURIComponent(nickname)}`);
    const data = await res.json();
    userHistory.value = data.apostas || [];
    isModalOpen.value = true;
  } catch(e) { alert('Erro: ' + e.message); }
};

const getHistoryBySeason = () => {
  const grouped = {};
  userHistory.value.forEach(a => {
    if (!grouped[a.season]) grouped[a.season] = [];
    grouped[a.season].push(a);
  });
  return grouped;
};

const reuseNumbers = (nums) => {
  form.value.fixos = nums.join(' ');
  form.value.qtd = nums.length;
  updateCost();
};

onMounted(() => {
  fetchData();
  intervalId = setInterval(fetchData, 5000);
});

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId);
});
</script>
