(function() {
	function A(b) {
		b = [].concat(a.values(b.proPlayers)).concat(a.values(b.proTeams));
		a.each(b, function(b) {
			var l = [b.statsBySeason].concat(a.values(b.statsByWeek)).concat(a.values(b.statsByMatch)),
				q = void 0 === b.proTeamId ? f.constants.ALL_TEAM_STAT_NAMES : f.constants.ALL_PLAYER_STAT_NAMES;
			a.each(l, function(b) {
				b.livePoints = b.actualPoints;
				a.each(q, function(a) {
					a = b[a];
					a.liveValue = a.actualValue;
					a.livePoints = a.actualPoints
				})
			})
		})
	}

	function y(a, k) {
		B(a, k, {}, [], [])
	}

	function B(b, k, m, q, n) {
		var p = function(b) {
			return a.chain(b).omit("RESERVE").values().flatten().map(function(a) {
				a =
					f.pros.getPro(k, a.targetType, a.targetId);
				return void 0 === a.proTeamId ? a : f.pros.getPro(k, "team", a.proTeamId)
			}).value()
		};
		l(".stat-listener-fantasy-team-total").each(function(k, f) {
			var d = l(f),
				c = d.data("week"),
				g = d.data("fantasyMatchId"),
				e = d.data("fantasyTeamId"),
				u = m[c],
				c = b.fantasyMatches[g];
			if (void 0 !== u && void 0 !== c) {
				if (c.leftTeam.id === e) e = c.leftTeam;
				else if (c.rightTeam.id === e) e = c.rightTeam;
				else return;
				c = a.chain(e.roster).omit("RESERVE").values().flatten().value();
				e = a.chain(c).filter(function(a) {
					return "team" ===
						a.targetType
				}).map(function(a) {
					return a.targetId
				}).value();
				c = a.chain(c).filter(function(a) {
					return "player" === a.targetType
				}).map(function(a) {
					return a.targetId
				}).value();
				e = a.reduce(e, function(c, e) {
					var d = u.byTeam[e];
					return void 0 === d ? c : c + a.reduce(b.pointsPerStat, function(a, b, c) {
						b *= d[c];
						return isNaN(b) ? a : a + b
					}, 0)
				}, 0);
				c = a.reduce(c, function(c, e) {
					var d = u.byPlayer[e];
					return void 0 === d ? c : c + a.reduce(b.pointsPerStat, function(a, b, c) {
						b *= d[c];
						return isNaN(b) ? a : a + b
					}, 0)
				}, 0);
				x(d.find(".formatted-points"), e + c)
			}
		});
		l(".stat-listener-pro-team-total").each(function(k,
			f) {
			var d = l(f),
				c = d.data("week"),
				g = d.data("proTeamId"),
				c = m[c];
			if (void 0 !== c) {
				var e = c.byTeam[g];
				void 0 !== e && (g = a.reduce(b.pointsPerStat, function(a, b, c) {
					b *= e[c];
					return isNaN(b) ? a : a + b
				}, 0), x(d.find(".formatted-points"), g))
			}
		});
		l(".stat-listener-pro-player-total").each(function(k, f) {
			var d = l(f),
				c = d.data("week"),
				g = d.data("proPlayerId"),
				c = m[c];
			if (void 0 !== c) {
				var e = c.byPlayer[g];
				void 0 !== e && (g = a.reduce(b.pointsPerStat, function(a, b, c) {
					b *= e[c];
					return isNaN(b) ? a : a + b
				}, 0), x(d.find(".formatted-points"), g))
			}
		});
		l(".stat-listener-fantasy-team-matches-remaining").each(function(k,
			f) {
			var d = l(f),
				c = d.data("week"),
				g = d.data("fantasyMatchId"),
				e = d.data("fantasyTeamId"),
				n = d.data("totalMatches"),
				q = d.data("completedMatches"),
				c = m[c],
				g = b.fantasyMatches[g];
			if (void 0 !== c && void 0 !== g) {
				var h;
				g.leftTeam.id === e ? h = g.leftTeam : g.rightTeam.id === e && (h = g.rightTeam);
				if (void 0 !== h) {
					var e = p(h.roster),
						t = a.map(c.byTeam, function(a, b) {
							return +b
						}),
						e = a.reduce(e, function(b, c) {
							return a.contains(t, c.id) ? b + 1 : b
						}, 0);
					d.text(n - q - e)
				}
			}
		});
		l(".stat-listener-pro-match-in-progress").each(function(b, h) {
			var d = l(h),
				c = d.data("week"),
				g = d.data("proTeamId");
			if (void 0 !== f.pros.getPro(k, "team", g)) {
				var e = a.filter(k.proMatches, function(b) {
						return b.week === c && a.contains([b.redTeamId, b.blueTeamId], g)
					}),
					m = a.map(e, function(a) {
						return a.id
					}),
					p = a.reduce(e, function(a, b) {
						return b.complete ? a + 1 : a
					}, 0),
					w = a.intersection(m, n).length,
					t = a.intersection(m, a.difference(q, n)).length;
				d.find(".completed-dot").each(function(a, b) {
					var c = l(b);
					a < p + w ? (c.removeClass("pending"), c.removeClass("in-progress"), c.hasClass("complete") || c.addClass("complete")) : a < p + w + t ? (c.removeClass("pending"),
						c.removeClass("complete"), c.hasClass("in-progress") || c.addClass("in-progress")) : (c.removeClass("in-progress"), c.removeClass("complete"), c.hasClass("pending") || c.addClass("pending"))
				})
			}
		})
	}
	var C = window.Math,
		a = window.Underscore,
		l = window.$,
		f = window.fantasy;
	f.livestats = {
		resetLiveStats: A,
		beginLiveStats: function(b, k) {
			if (null === h) {
				var m = l(".fantasy-main"),
					q = m.find(".fantasy-matchups"),
					n = m.data("statsWebsocketEnabled"),
					p = m.data("statsWebsocketUri");
				if (n && !(0 === b.currentWeek || b.currentWeek > k.numberOfWeeks)) {
					var r = {},
						z = {},
						d = a.indexBy(a.values(k.proMatches), "riotId"),
						c = a.indexBy(a.values(k.proPlayers), "riotId"),
						g = {},
						e = [],
						u = [],
						x = function(b) {
							b.complete || -1 !== a.indexOf(e, b.id) || e.push(b.id);
							var c = g[b.id];
							void 0 === c && (c = {
								byPlayer: {},
								byTeam: {}
							}, g[b.id] = c);
							return c
						},
						w = function(c, e) {
							a.each(e, function(a, d) {
								var k = L[d];
								void 0 !== k && (d = k);
								isNaN(b.pointsPerStat[d]) || (c[d] = a, ("kills" === d || "assists" === d) && 10 <= a && (c.killOrAssistBonus = 1), "gameLength" === d && 1800 > a && 0 < a && (0 < e.matchVictory || 0 < c.matchVictory && 0 !== e.matchVictory) &&
									(c.quickWinBonus = 1))
							})
						},
						t = 0,
						y = function(a) {
							var b = [1, 1, 2, 3, 5, 8, 13, 21, 34, 55];
							return C.ceil((1 - .2 * C.random() + .1) * b[a >= b.length ? b.length - 1 : a] * 1E3)
						},
						J = function(b) {
							a.each(b, function(b, e) {
								if (null !== b && "" !== b) {
									a.has(b, "matchId") && (r[e] = b.matchId);
									var g = d[r[e]];
									if (g) {
										var f = x(g);
										a.each(b.playerStats, function(a, b) {
											a.playerId && (z[e] || (z[e] = {}), z[e][b] = a.playerId);
											var d = c[z[e][b]];
											if (d) {
												var v = k.proTeams[d.proTeamId];
												void 0 !== v.statsByMatch[g.id] && 0 < v.statsByMatch[g.id].matchesPlayed.actualValue || (v = f.byPlayer[d.id],
													void 0 === v && (v = D(), f.byPlayer[d.id] = v), w(v, a))
											}
										});
										a.each(b.teamStats, function(a, b) {
											var c = 200 == b ? k.proTeams[g.redTeamId] : k.proTeams[g.blueTeamId];
											if (null != c && !(void 0 !== c.statsByMatch[g.id] && 0 < c.statsByMatch[g.id].matchesPlayed.actualValue)) {
												var d = f.byTeam[c.id];
												void 0 === d && (d = E(), f.byTeam[c.id] = d);
												w(d, a)
											}
										});
										b.gameComplete && !g.complete && (u = a.union(u, [g.id]))
									}
								}
							})
						},
						K = function(c) {
							a.each(c, function(c, d) {
								var e = l(".fantasy-main").find(".fantasy-matchups").find(".week" + d);
								0 !== e.length && e.data("matchesStarted") &&
									(a.each(c.byPlayer, function(c, g) {
										var f = e.find(".pro-player-" + g);
										if (0 !== f.length) {
											var l = k.proPlayers[g];
											f.data("fantasyTeam");
											var m = 0;
											a.each(b.pointsPerStat, function(a, b) {
												var d = c[b];
												isNaN(d) || (m += d * a)
											});
											var h = l.statsByWeek[d];
											h.livePoints += m;
											a.each(F, function(a) {
												var d = h[a],
													e = c[a];
												a = b.pointsPerStat[a];
												d.liveValue += e;
												d.livePoints += a * e
											})
										}
									}), a.each(c.byTeam, function(c, g) {
										var f = e.find(".pro-team-" + g);
										if (0 !== f.length) {
											var m = k.proTeams[g];
											f.data("fantasyTeam");
											var l = 0;
											a.each(b.pointsPerStat, function(a, b) {
												var d =
													c[b];
												isNaN(d) || (l += d * a)
											});
											var h = m.statsByWeek[d];
											h.livePoints += l;
											a.each(G, function(a) {
												var d = h[a],
													e = c[a];
												a = b.pointsPerStat[a];
												d.liveValue += e;
												d.livePoints += a * e
											})
										}
									}))
							})
						},
						H = function(a) {
							t = 1;
							try {
								J(JSON.parse(a.data))
							} catch (c) {
								return
							}
							A(k);
							var d = M(k, g);
							K(d);
							f.livestats.updateStatListeners = function(a, b) {
								B(a, b, d, e, u)
							};
							f.livestats.updateStatListeners(b, k);
							q.find(".stats-table").trigger("refresh");
							f.updates.updateMatchupsPointsBars(b, k, b.currentWeek)
						},
						I = function(a) {
							!1 === a.wasClean && 24 > t ? (5 == t && (a = f.util.localize("ui.misc.websocket.connect.warn"),
								f.util.alertText("info", a)), f.timeouts.clear("statsWebsocketReconnect"), f.timeouts.set("statsWebsocketReconnect", function() {
								h = new WebSocket(p);
								h.onclose = I;
								h.onmessage = H;
								t++
							}, y(t))) : !1 === a.wasClean && 24 == t && (a = f.util.localize("ui.misc.websocket.connect.error"), f.util.alertText("error", a))
						};
					h = new WebSocket(p);
					h.onclose = I;
					h.onmessage = H
				}
			}
		},
		stopLiveStats: function() {
			f.livestats.updateStatListeners = y;
			null !== h && (h.close(), h = null)
		},
		updateStatListeners: y
	};
	var G = "firstBlood towerKills baronKills dragonKills matchVictory quickWinBonus".split(" "),
		F = "kills deaths assists minionKills tripleKills quadraKills pentaKills killOrAssistBonus".split(" "),
		L = {
			mk: "minionKills",
			towersKilled: "towerKills",
			baronsKilled: "baronKills",
			dragonsKilled: "dragonKills"
		},
		x = function(a, k) {
			var h = a.data("number"),
				q = a.data("liveNumber"),
				n = h + k;
			if (q !== n) {
				var p = a.data("pointsType");
				a.data("liveNumber", n);
				a.html(l(f.util.formatPoints(n, p)).html());
				l({
					value: q
				}).animate({
					value: n
				}, {
					duration: 500,
					easing: "easeInOutCubic",
					step: function() {
						a.html(l(f.util.formatPoints(this.value, p)).html())
					},
					done: function() {
						a.html(l(f.util.formatPoints(n, p)).html())
					}
				})
			}
		},
		D = function() {
			var b = {};
			a.each(F, function(a) {
				b[a] = 0
			});
			return b
		},
		E = function() {
			var b = {};
			a.each(G, function(a) {
				b[a] = 0
			});
			return b
		},
		M = function(b, f) {
			var h = {};
			a.each(f, function(f, k) {
				var l = b.proMatches[k].week,
					r = h[l];
				void 0 === r && (r = {
					byPlayer: {},
					byTeam: {}
				}, h[l] = r);
				a.each(f.byPlayer, function(b, d) {
					var c = r.byPlayer[d];
					void 0 === c && (c = D(), r.byPlayer[d] = c);
					a.each(b, function(a, b) {
						c[b] += a
					})
				});
				a.each(f.byTeam, function(b, d) {
					var c = r.byTeam[d];
					void 0 === c && (c =
						E(), r.byTeam[d] = c);
					a.each(b, function(a, b) {
						c[b] += a
					})
				})
			});
			return h
		},
		h = null
})();
