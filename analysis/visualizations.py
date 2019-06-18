import networkx as nx
import numpy as np
import matplotlib.pyplot as plt

G = nx.Graph()

G.add_edge('NCAA Champs', 'West', weight=-23)
G.add_edge('NCAA Champs', 'South', weight=124)
G.add_edge('NCAA Champs', 'Southeast', weight=-10)
G.add_edge('NCAA Champs', 'South Central', weight=43)
G.add_edge('NCAA Champs', 'Northeast', weight=128)
G.add_edge('NCAA Champs', 'Mountain', weight=-31)
G.add_edge('NCAA Champs', 'Mid-Atlantic', weight=23)
G.add_edge('NCAA Champs', 'Great Lakes', weight=42)

elarge = [(u, v) for (u, v, d) in G.edges(data=True) if d['weight'] > 0]
esmall = [(u, v) for (u, v, d) in G.edges(data=True) if d['weight'] <= 0]

pos = nx.spring_layout(G)

nx.draw_networkx_nodes(G, pos, node_size=700)

nx.draw_networkx_edges(G, pos, edgelist=elarge, width=6)
nx.draw_networkx_edges(G, pos, edgelist=esmall, width=6, alpha=0.5, edge_color='b', style='dashed')

nx.draw_networkx_labels(G, pos, font_size=20, font_family='sans-serif')

plt.axis('off')
plt.show()