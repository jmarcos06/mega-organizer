        // Global variables
        let currentSeason = '';
        let allSeasons = [];

        // --- Cost calculation ---
        const qtdInput = document.getElementById('qtd');
        const display = document.getElementById('costDisplay');

        async function atualizarCusto() {
            const val = parseInt(qtdInput.value);
            if (val < 6 || val > 20) {
                display.innerText = "Muitos $$$";
                return;
            }
            try {
                const res = await fetch('/api/custo?qtd=' + val);
                const data = await res.json();
                display.innerText = 'R$ ' + data.custo.toLocaleString('pt-BR', {minimumFractionDigits: 2});
            } catch (err) {
                display.innerText = "Erro...";
                console.error(err);
            }
        }

        qtdInput.addEventListener('input', atualizarCusto);
        atualizarCusto(); // Initialize

        // --- Form submission ---
        document.getElementById('betForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const btn = e.target.querySelector('button');
            const originalText = btn.innerHTML;
            btn.disabled = true;
            btn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Processando...';

            const formData = new FormData(e.target);
            const data = Object.fromEntries(formData.entries());

            try {
                const res = await fetch('/api/apostar', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify(data)
                });
                
                const json = await res.json();
                
                if (!res.ok) throw new Error(json.error || 'Erro desconhecido');

                if (json.colisao) {
                    alert('⚠️ ATENÇÃO: O usuário ' + json.colisao_com + ' já fez esse jogo exato!');
                } else {
                    confetti({ particleCount: 100, spread: 70, origin: { y: 0.6 } });
                }

                e.target.reset();
                // Reset season to default
                document.getElementById('season').value = 'Mega da Virada 2025';
                atualizarCusto(); // Reset cost display
                await fetchData(); // Refresh data

            } catch (err) {
                alert('Erro: ' + err.message);
            } finally {
                btn.disabled = false;
                btn.innerHTML = originalText;
            }
        });

        // --- Season filter ---
        document.getElementById('seasonFilter').addEventListener('change', async (e) => {
            currentSeason = e.target.value;
            await fetchData();
        });

        // --- Data fetching ---
        async function fetchData() {
            try {
                const seasonParam = currentSeason ? `?season=${encodeURIComponent(currentSeason)}` : '';
                const res = await fetch('/api/dados' + seasonParam);
                const data = await res.json();
                render(data);
            } catch (err) {
                console.error("Erro ao buscar dados", err);
            }
        }

        function render(data) {
            // Update seasons in filters if changed
            if (JSON.stringify(data.seasons) !== JSON.stringify(allSeasons)) {
                allSeasons = data.seasons;
                updateSeasonSelectors(data.seasons);
            }

            // Stats
            document.getElementById('totalMoney').innerText = data.total_gasto.toLocaleString('pt-BR', {style: 'currency', currency: 'BRL'});
            document.getElementById('totalGames').innerText = data.total_jogos;

            // Hot Numbers
            const hotDiv = document.getElementById('hotNumbers');
            hotDiv.innerHTML = data.numeros_quentes.map(n => 
                '<span class="bg-yellow-600 text-yellow-100 text-xs px-2 py-1 rounded-full font-bold">' + n.numero + ' (' + n.qtd + 'x)</span>'
            ).join('');

            // Feed
            renderFeed(data.ultimas_apostas);
        }

        function updateSeasonSelectors(seasons) {
            // Update season filter
            const seasonFilter = document.getElementById('seasonFilter');
            const currentValue = seasonFilter.value;
            seasonFilter.innerHTML = '<option value="">Todas as temporadas</option>';
            seasons.forEach(season => {
                const option = document.createElement('option');
                option.value = season;
                option.textContent = season;
                seasonFilter.appendChild(option);
            });
            seasonFilter.value = currentValue;

            // Update form season selector
            const seasonForm = document.getElementById('season');
            const formValue = seasonForm.value;
            seasons.forEach(season => {
                if (!Array.from(seasonForm.options).find(opt => opt.value === season)) {
                    const option = document.createElement('option');
                    option.value = season;
                    option.textContent = season;
                    seasonForm.appendChild(option);
                }
            });
            seasonForm.value = formValue;
        }

        function renderFeed(apostas) {
            const feedDiv = document.getElementById('feed');
            if (apostas.length === 0) {
                feedDiv.innerHTML = '<div class="text-center text-slate-500">Nenhuma aposta ainda. Seja o primeiro!</div>';
                return;
            }

            feedDiv.innerHTML = apostas.map(aposta => {
                const numsHtml = aposta.numeros.map(n => 
                    '<span class="inline-block w-8 h-8 leading-8 text-center rounded-full bg-slate-700 text-slate-200 font-bold text-sm border border-slate-600 shadow-sm">' + n + '</span>'
                ).join('');

                const tipoBadge = aposta.tipo === 'Simples' 
                    ? '<span class="text-xs bg-green-900 text-green-300 px-2 py-0.5 rounded border border-green-700">Simples</span>'
                    : '<span class="text-xs bg-purple-900 text-purple-300 px-2 py-0.5 rounded border border-purple-700">Desdobramento</span>';

                const seasonBadge = '<span class="text-xs bg-blue-900 text-blue-300 px-1 py-0.5 rounded border border-blue-700">' + aposta.season + '</span>';

                const deleteBtn = '<button onclick="deleteAposta(\'' + aposta.id + '\')" class="text-red-400 hover:text-red-300 text-xs ml-2" title="Deletar aposta"><i class="fas fa-trash"></i></button>';

                return '<div class="bg-slate-750 hover:bg-slate-700 transition p-4 rounded-lg border border-slate-700/50">' +
                    '<div class="flex justify-between items-start mb-2">' +
                        '<div>' +
                            '<span class="font-bold text-white text-lg mr-2 cursor-pointer hover:text-green-400" onclick="showUserHistory(\'' + aposta.nickname + '\')">' + aposta.nickname + '</span>' +
                            tipoBadge + ' ' + seasonBadge +
                        '</div>' +
                        '<div class="text-right text-xs text-slate-400 flex items-center">' +
                            '<div>' +
                                '<div>' + new Date(aposta.data).toLocaleTimeString() + '</div>' +
                                '<div class="text-slate-500 text-[10px]">' + aposta.id + '</div>' +
                            '</div>' + deleteBtn +
                        '</div>' +
                    '</div>' +
                    '<div class="flex flex-wrap gap-2 mb-2">' + numsHtml + '</div>' +
                    '<div class="text-xs text-slate-500 flex justify-between">' +
                        '<span>Custo: ' + aposta.custo.toLocaleString('pt-BR', {style: 'currency', currency: 'BRL'}) + '</span>' +
                        '<button onclick="reuseNumbers([' + aposta.numeros.join(',') + '])" class="text-blue-400 hover:text-blue-300"><i class="fas fa-copy"></i> Reusar números</button>' +
                    '</div>' +
                '</div>';
            }).join('');
        }

        // --- User History Modal ---
        async function showUserHistory(nickname) {
            try {
                const res = await fetch(`/api/usuario/historico?nickname=${encodeURIComponent(nickname)}`);
                const data = await res.json();
                
                document.getElementById('modalTitle').textContent = `Histórico de ${nickname}`;
                
                if (data.apostas.length === 0) {
                    document.getElementById('modalContent').innerHTML = 
                        '<div class="text-center text-slate-500">Este usuário ainda não fez apostas.</div>';
                } else {
                    const groupedBySeason = {};
                    data.apostas.forEach(aposta => {
                        if (!groupedBySeason[aposta.season]) {
                            groupedBySeason[aposta.season] = [];
                        }
                        groupedBySeason[aposta.season].push(aposta);
                    });

                    let html = '';
                    Object.keys(groupedBySeason).sort().forEach(season => {
                        html += `<div class="mb-6">
                            <h3 class="text-lg font-bold text-green-400 mb-3 border-b border-slate-600 pb-2">${season}</h3>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">`;
                        
                        groupedBySeason[season].forEach(aposta => {
                            const numsHtml = aposta.numeros.map(n => 
                                '<span class="inline-block w-6 h-6 leading-6 text-center rounded-full bg-slate-600 text-slate-200 font-bold text-xs border border-slate-500">' + n + '</span>'
                            ).join(' ');
                            
                            html += `<div class="bg-slate-700 p-3 rounded border border-slate-600">
                                <div class="flex justify-between items-center mb-2">
                                    <span class="text-xs text-slate-400">${new Date(aposta.data).toLocaleString()}</span>
                                    <span class="text-xs ${aposta.tipo === 'Simples' ? 'text-green-400' : 'text-purple-400'}">${aposta.tipo}</span>
                                </div>
                                <div class="mb-2">${numsHtml}</div>
                                <div class="flex justify-between items-center text-xs">
                                    <span class="text-slate-500">R$ ${aposta.custo.toFixed(2)}</span>
                                    <button onclick="reuseNumbers([${aposta.numeros.join(',')}]); closeUserModal();" class="text-blue-400 hover:text-blue-300">
                                        <i class="fas fa-copy"></i> Reusar
                                    </button>
                                </div>
                            </div>`;
                        });
                        html += '</div></div>';
                    });

                    document.getElementById('modalContent').innerHTML = html;
                }
                
                document.getElementById('userModal').classList.remove('hidden');
            } catch (err) {
                alert('Erro ao carregar histórico: ' + err.message);
            }
        }

        function closeUserModal() {
            document.getElementById('userModal').classList.add('hidden');
        }

        // --- Utility functions ---
        function reuseNumbers(numbers) {
            document.getElementById('fixos').value = numbers.join(' ');
            document.getElementById('qtd').value = numbers.length;
            atualizarCusto();
        }

        async function deleteAposta(id) {
            if (!confirm('Tem certeza que deseja deletar esta aposta?')) return;
            
            const token = prompt('Digite o token de acesso para deletar:');
            if (!token) return;

            try {
                const res = await fetch(`/api/aposta/deletar?id=${id}&token=${token}`, {
                    method: 'DELETE'
                });
                
                if (!res.ok) {
                    const error = await res.json();
                    throw new Error(error.error || 'Erro ao deletar');
                }
                
                await fetchData(); // Refresh data
                alert('Aposta deletada com sucesso!');
            } catch (err) {
                alert('Erro ao deletar: ' + err.message);
            }
        }

        // Close modal when clicking outside
        document.getElementById('userModal').addEventListener('click', (e) => {
            if (e.target.id === 'userModal') {
                closeUserModal();
            }
        });

        // --- Initialize ---
        setInterval(fetchData, 5000); // Auto-refresh every 5 seconds
        fetchData(); // Initial load